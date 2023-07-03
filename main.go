package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	"github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// if err := godotenv.Load(".env"); err != nil {
	// 	panic(fmt.Sprintf("Failed env : %v", err))
	// }
	// dsn := os.Getenv("DB")
	dsn := "postgresql://doadmin:AVNS_6LM8d5bO2GCwNym-Hkl@indoora-db-do-user-13831734-0.b.db.ondigitalocean.com:25060/defaultdb"
	connstring := os.ExpandEnv(dsn)
	database, err := gorm.Open(postgres.Open(connstring), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}
	fmt.Println("successfully connected")

	database.AutoMigrate(&VibraStore{})
	DB = database
}

type handler struct {
	json bool
}

type Axis struct {
	SenEvent     int64 `json:"SenEvent"`
	OAVelocity   int64 `json:"OAVelocity"`
	Peakmg       int64 `json:"Peakmg"`
	RMSmg        int64 `json:"RMSmg"`
	Kurtosis     int64 `json:"Kurtosis"`
	CrestFactor  int64 `json:"CrestFactor"`
	Skewness     int64 `json:"Skewness"`
	Deviation    int64 `json:"Deviation"`
	Displacement int64 `json:"Peak-to-Peak Displacement"`
}

type TempHumi struct {
	Range  int64 `json:"Range"`
	Status int64 `json:"Status"`
	Event  int64 `json:"Event"`
	SenVal int64 `json:"SenVal"`
}

type Device struct {
	Events      int64 `json:"Events"`
	PowerSrc    int64 `json:"PowerSrc"`
	BatteryVolt int64 `json:"BatteryVold"`
	Time        int64 `json:"Time"`
}

type Accelerometer struct {
	XAxis    Axis  `json:"X-Axis"`
	YAxis    Axis  `json:"Y-Axis"`
	ZAxis    Axis  `json:"Z-Axis"`
	LogIndex int64 `json:"LogIndex"`
	Time     int64 `json:"Time"`
}

type Vibra struct {
	SequenceNumber int64         `json:"SequenceNumber"`
	TotalLength    int64         `json:"TotalLength"`
	SourceAddress  int64         `json:"SourceAddress"`
	TempHumi       TempHumi      `json:"TempHumi"`
	Accelerometer  Accelerometer `json:"Accelerometer"`
	Device         Device        `json:"Device"`
}

type VibraStore struct {
	ID                int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	ApplicationID     uint64 `json:"applicationID"`
	ApplicationName   string `json:"applicationName"`
	DeviceName        string `json:"deviceName"`
	DeviceProfileName string `json:"deviceProfileName"`
	DeviceProfileID   string `json:"deviceProfileID"`
	DevEUI            []byte `json:"devEUI"`
	Data              string `json:"data"`
	SequenceNumber    int64  `json:"sequence_number"`
	TotalLength       int64  `json:"total_length"`
	SourceAddress     int64  `json:"source_address"`
	TempHumiRange     int64  `json:"temp_humi_range"`
	TempHumiStatus    int64  `json:"temp_humi_status"`
	TempHumiEvent     int64  `json:"temp_humi_event"`
	TempHumiSenVal    int64  `json:"temp_humi_sen_val"`
	XAxisSenEvent     int64  `json:"x_axis_sen_event"`
	XAxisOAVelocity   int64  `json:"x_axis_oavelocity"`
	XAxisPeakmg       int64  `json:"x_axis_peakmg"`
	XAxisRMSmg        int64  `json:"x_axis_rmsmg"`
	XAxisKurtosis     int64  `json:"x_axis_kurtosis"`
	XAxisCrestFactor  int64  `json:"x_axis_crest_factor"`
	XAxisSkewness     int64  `json:"x_axis_skewness"`
	XAxisDeviation    int64  `json:"x_axis_deviation"`
	XAxisDisplacement int64  `json:"x_axis_displacement"`

	YAxisSenEvent     int64 `json:"y_axis_sen_event"`
	YAxisOAVelocity   int64 `json:"y_axis_oavelocity"`
	YAxisPeakmg       int64 `json:"y_axis_peakmg"`
	YAxisRMSmg        int64 `json:"y_axis_rmsmg"`
	YAxisKurtosis     int64 `json:"y_axis_kurtosis"`
	YAxisCrestFactor  int64 `json:"y_axis_crest_factor"`
	YAxisSkewness     int64 `json:"y_axis_skewness"`
	YAxisDeviation    int64 `json:"y_axis_deviation"`
	YAxisDisplacement int64 `json:"y_axis_displacement"`

	ZAxisSenEvent     int64 `json:"z_axis_sen_event"`
	ZAxisOAVelocity   int64 `json:"z_axis_oavelocity"`
	ZAxisPeakmg       int64 `json:"z_axis_peakmg"`
	ZAxisRMSmg        int64 `json:"z_axis_rmsmg"`
	ZAxisKurtosis     int64 `json:"z_axis_kurtosis"`
	ZAxisCrestFactor  int64 `json:"z_axis_crest_factor"`
	ZAxisSkewness     int64 `json:"z_axis_skewness"`
	ZAxisDeviation    int64 `json:"z_axis_deviation"`
	ZAxisDisplacement int64 `json:"z_axis_displacement"`

	LogIndex int64 `json:"log_index"`
	Time1    int64 `json:"time1"`

	DeviceEvents      int64     `json:"device_events"`
	DevicePowerSrc    int64     `json:"device_power_src"`
	DeviceBatteryVolt int64     `json:"device_battery_volt"`
	Time2             int64     `json:"time2"`
	CreatedAt         time.Time `json:"created_at"`
}

func (c *VibraStore) TableName() string {
	return "vibra_sensor_erdenet"
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	event := r.URL.Query().Get("event")

	switch event {
	case "up":
		err = h.up(b)
	case "join":
		err = h.join(b)
	default:
		log.Printf("handler for event %s is not implemented", event)
		return
	}

	if err != nil {
		log.Printf("handling event '%s' returned error: %s", event, err)
	}
}

func (h *handler) up(b []byte) error {
	var up integration.UplinkEvent
	if err := h.unmarshal(b, &up); err != nil {
		return err
	}
	fmt.Printf("Uplink received from %t with payload: %s\n", up.ConfirmedUplink, hex.EncodeToString(up.Data))

	fmt.Println("HEX :: ", hex.EncodeToString(up.Data))

	cmd := exec.Command("node", "/root/vibra-erdenet/wise_engine.js", hex.EncodeToString(up.Data))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error executing Node.js script: %s", err)
		log.Printf("Standard error: %s", stderr.String())
		return err
	}

	otpu := stdout.String()
	log.Printf("Output: %s", otpu)

	output, _ := cmd.Output()
	fmt.Println("output : ", output)
	var vibra Vibra
	var store VibraStore
	if err := json.Unmarshal(output, &vibra); err != nil {
		fmt.Println("Error:", err.Error())

	} else {

		store.SequenceNumber = vibra.SequenceNumber
		store.TotalLength = vibra.TotalLength
		store.SourceAddress = vibra.SourceAddress

		store.TempHumiRange = vibra.TempHumi.Range
		store.TempHumiEvent = vibra.TempHumi.Event
		store.TempHumiStatus = vibra.TempHumi.Status
		store.TempHumiSenVal = vibra.TempHumi.SenVal

		store.XAxisSenEvent = vibra.Accelerometer.XAxis.SenEvent
		store.XAxisOAVelocity = vibra.Accelerometer.XAxis.OAVelocity
		store.XAxisPeakmg = vibra.Accelerometer.XAxis.Peakmg
		store.XAxisRMSmg = vibra.Accelerometer.XAxis.RMSmg
		store.XAxisKurtosis = vibra.Accelerometer.XAxis.Kurtosis
		store.XAxisCrestFactor = vibra.Accelerometer.XAxis.CrestFactor
		store.XAxisSkewness = vibra.Accelerometer.XAxis.Skewness
		store.XAxisDeviation = vibra.Accelerometer.XAxis.Deviation
		store.XAxisDisplacement = vibra.Accelerometer.XAxis.Displacement

		store.YAxisSenEvent = vibra.Accelerometer.YAxis.SenEvent
		store.YAxisOAVelocity = vibra.Accelerometer.YAxis.OAVelocity
		store.YAxisPeakmg = vibra.Accelerometer.YAxis.Peakmg
		store.YAxisRMSmg = vibra.Accelerometer.YAxis.RMSmg
		store.YAxisKurtosis = vibra.Accelerometer.YAxis.Kurtosis
		store.YAxisCrestFactor = vibra.Accelerometer.YAxis.CrestFactor
		store.YAxisSkewness = vibra.Accelerometer.YAxis.Skewness
		store.YAxisDeviation = vibra.Accelerometer.YAxis.Deviation
		store.YAxisDisplacement = vibra.Accelerometer.YAxis.Displacement

		store.ZAxisSenEvent = vibra.Accelerometer.ZAxis.SenEvent
		store.ZAxisOAVelocity = vibra.Accelerometer.ZAxis.OAVelocity
		store.ZAxisPeakmg = vibra.Accelerometer.ZAxis.Peakmg
		store.ZAxisRMSmg = vibra.Accelerometer.ZAxis.RMSmg
		store.ZAxisKurtosis = vibra.Accelerometer.ZAxis.Kurtosis
		store.ZAxisCrestFactor = vibra.Accelerometer.ZAxis.CrestFactor
		store.ZAxisSkewness = vibra.Accelerometer.ZAxis.Skewness
		store.ZAxisDeviation = vibra.Accelerometer.ZAxis.Deviation
		store.ZAxisDisplacement = vibra.Accelerometer.ZAxis.Displacement

		store.LogIndex = vibra.Accelerometer.LogIndex
		store.Time1 = vibra.Accelerometer.Time

		store.DeviceEvents = vibra.Device.Events
		store.DevicePowerSrc = vibra.Device.PowerSrc
		store.DeviceBatteryVolt = vibra.Device.BatteryVolt
		store.Time2 = vibra.Device.Time

		store.ApplicationID = (up.ApplicationId)
		store.ApplicationName = up.ApplicationName
		store.DeviceName = up.DeviceName
		store.DeviceProfileID = up.DeviceProfileId
		store.DeviceProfileName = up.DeviceProfileName
		store.DevEUI = (up.DevEui)
		store.Data = hex.EncodeToString(up.Data)
		store.CreatedAt = time.Now()

		if err := DB.Create(&store).Error; err != nil {
			fmt.Println("error :", err.Error())
		}
	}
	return nil
}

func (h *handler) join(b []byte) error {
	var join integration.JoinEvent
	if err := h.unmarshal(b, &join); err != nil {
		return err
	}
	fmt.Printf("Device %s joined with DevAddr %s\n", hex.EncodeToString(join.DevEui), hex.EncodeToString(join.DevAddr))
	return nil
}

func (h *handler) unmarshal(b []byte, v proto.Message) error {
	if h.json {
		unmarshaler := &jsonpb.Unmarshaler{
			AllowUnknownFields: true, // we don't want to fail on unknown fields
		}
		return unmarshaler.Unmarshal(bytes.NewReader(b), v)
	}
	return proto.Unmarshal(b, v)
}

func main() {
	ConnectDatabase()
	// json: false   - to handle Protobuf payloads (binary)
	// json: true    - to handle JSON payloads (Protobuf JSON mapping)
	http.Handle("/", &handler{json: true})
	log.Fatal(http.ListenAndServe(":8090", nil))
}
