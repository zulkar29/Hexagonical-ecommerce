import React from 'react';
import { useParams } from 'react-router-dom';

export default function SupportDetail() {
  const { id } = useParams();
  return <div>Support Ticket Details Page: {id}</div>;
} 