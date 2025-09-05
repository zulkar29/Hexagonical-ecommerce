import React from 'react';
import { useParams } from 'react-router-dom';

export default function PaymentDetail() {
  const { id } = useParams();
  return <div>Payment Details Page: {id}</div>;
} 