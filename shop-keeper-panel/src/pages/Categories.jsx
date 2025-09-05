import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Plus, Edit, Trash2, Tags } from 'lucide-react';

const mockCategories = [
  { id: 1, name: 'Rice & Grains', description: 'All types of rice, lentils, and grains.' },
  { id: 2, name: 'Vegetables', description: 'Fresh vegetables and greens.' },
  { id: 3, name: 'Dairy', description: 'Milk, cheese, yogurt, and more.' },
  { id: 4, name: 'Cooking', description: 'Oils, spices, and essentials.' },
  { id: 5, name: 'Beverages', description: 'Tea, coffee, and drinks.' },
];

const Categories = () => {
  const [search, setSearch] = useState('');
  const filtered = mockCategories.filter(c =>
    c.name.toLowerCase().includes(search.toLowerCase()) ||
    c.description.toLowerCase().includes(search.toLowerCase())
  );
  return (
    <div className="space-y-6 p-6 bg-background min-h-screen">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Categories</h1>
          <p className="text-muted-foreground mt-1">Organize your products by category</p>
        </div>
        <Button>
          <Plus className="h-4 w-4 mr-2" /> Add Category
        </Button>
      </div>
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center mb-4">
            <Input
              placeholder="Search categories..."
              value={search}
              onChange={e => setSearch(e.target.value)}
              className="max-w-xs"
            />
          </div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Category</TableHead>
                <TableHead>Description</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map(category => (
                <TableRow key={category.id}>
                  <TableCell>
                    <div className="flex items-center space-x-3">
                      <Badge variant="secondary"><Tags className="h-4 w-4 mr-1" />{category.name}</Badge>
                    </div>
                  </TableCell>
                  <TableCell>{category.description}</TableCell>
                  <TableCell>
                    <Button size="sm" variant="outline" className="mr-2"><Edit className="h-4 w-4" /></Button>
                    <Button size="sm" variant="outline" className="text-destructive"><Trash2 className="h-4 w-4" /></Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
};

export default Categories; 