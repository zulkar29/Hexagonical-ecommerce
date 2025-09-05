import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Shield, Save, X, CheckSquare } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Checkbox } from '@/components/ui/checkbox';

const mockRoles = {
  1: { name: 'Admin', desc: 'Full access to all features', perms: ['View Users', 'Edit Users', 'Delete Users', 'Manage Settings'] },
  2: { name: 'Manager', desc: 'Manage users and tenants', perms: ['View Users', 'Edit Users', 'View Tenants', 'Edit Tenants'] },
  3: { name: 'Support', desc: 'Handle support tickets', perms: ['View Users', 'View Tenants'] },
  4: { name: 'Viewer', desc: 'Read-only access', perms: ['View Users', 'View Tenants', 'Access Analytics'] },
  5: { name: 'Suspended', desc: 'No access', perms: [] }
};

const mockPermissions = [
  'View Users',
  'Edit Users',
  'Delete Users',
  'View Tenants',
  'Edit Tenants',
  'Manage Billing',
  'Access Analytics',
  'Manage Settings'
];

export default function EditRole() {
  const { id } = useParams();
  const navigate = useNavigate();
  const role = mockRoles[id] || mockRoles[1];
  const [name, setName] = useState(role.name);
  const [desc, setDesc] = useState(role.desc);
  const [perms, setPerms] = useState(role.perms);

  const handlePermChange = (perm) => {
    setPerms((prev) => prev.includes(perm) ? prev.filter(p => p !== perm) : [...prev, perm]);
  };

  const handleSave = () => {
    // Simulate save
    navigate('/users/permissions');
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center gap-4 px-6">
          <Shield className="h-5 w-5 text-primary" />
          <h1 className="text-xl font-semibold">Edit Role</h1>
        </div>
      </div>
      <div className="flex-1 overflow-auto p-6 flex justify-center items-start">
        <Card className="w-full max-w-xl">
          <CardHeader>
            <CardTitle>Edit Role</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div>
              <label className="block text-sm font-medium mb-1">Role Name</label>
              <Input value={name} onChange={e => setName(e.target.value)} placeholder="Enter role name" />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Description</label>
              <Textarea value={desc} onChange={e => setDesc(e.target.value)} placeholder="Describe this role" />
            </div>
            <div>
              <label className="block text-sm font-medium mb-2">Permissions</label>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                {mockPermissions.map((perm) => (
                  <label key={perm} className="flex items-center gap-2">
                    <Checkbox checked={perms.includes(perm)} onCheckedChange={() => handlePermChange(perm)} />
                    <span>{perm}</span>
                  </label>
                ))}
              </div>
            </div>
            <div className="flex gap-2 justify-end">
              <Button variant="outline" onClick={() => navigate('/users/permissions')}><X className="h-4 w-4 mr-2" />Cancel</Button>
              <Button onClick={handleSave}><Save className="h-4 w-4 mr-2" />Save Changes</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
