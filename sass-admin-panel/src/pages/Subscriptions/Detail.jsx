import React from 'react';
import { useParams } from 'react-router-dom';

export default function SubscriptionDetail() {
  const { id } = useParams();
  return <div>Subscription Details Page: {id}</div>;
} 