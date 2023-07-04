import React from 'react';
import { Table } from 'antd';

const Tables = () => {
  const dataSource1 = [
    {
      key: '1',
      name: 'John Doe',
      age: 30,
      address: '123 Main St',
    },
    // Add more data as needed
  ];

  const dataSource2 = [
    {
      key: '1',
      company: 'ABC Corp',
      industry: 'IT',
      employees: 100,
    },
    // Add more data as needed
  ];

  const columns1 = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Age',
      dataIndex: 'age',
      key: 'age',
    },
    {
      title: 'Address',
      dataIndex: 'address',
      key: 'address',
    },
  ];

  const columns2 = [
    {
      title: 'Company',
      dataIndex: 'company',
      key: 'company',
    },
    {
      title: 'Industry',
      dataIndex: 'industry',
      key: 'industry',
    },
    {
      title: 'Employees',
      dataIndex: 'employees',
      key: 'employees',
    },
  ];

  return (
    <div>
      <Table dataSource={dataSource1} columns={columns1} />
      <Table dataSource={dataSource2} columns={columns2} />
    </div>
  );
};

export default Tables;
