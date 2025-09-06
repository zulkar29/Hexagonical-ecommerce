import { Users, ShieldCheck, Edit, Trash2, Plus, AlertCircle } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

const roles = [
  { id: 1, name: "Admin", description: "Full access to all features.", users: 2 },
  { id: 2, name: "Manager", description: "Manage products, orders, and customers.", users: 4 },
  { id: 3, name: "Support", description: "Handle support tickets and customer queries.", users: 3 },
  { id: 4, name: "Viewer", description: "Read-only access to dashboard.", users: 1 },
];

export default function UserRoles() {
  const navigate = useNavigate();

  const handleEditRole = (roleId) => {
    navigate(`/user-roles/${roleId}/edit`);
  };

  const handleDeleteRole = (roleId) => {
    // TODO: Add delete confirmation dialog
    console.log("Delete role:", roleId);
  };

  return (
    <div className="space-y-4 sm:space-y-6 p-3 sm:p-6">
      <Card>
        <CardHeader className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <CardTitle className="text-lg sm:text-xl">User Roles</CardTitle>
            <CardDescription className="text-sm">
              Manage access levels and permissions for your team members
            </CardDescription>
            <div className="flex items-center gap-2 text-sm text-amber-600 bg-amber-50 p-2 px-3 rounded-lg mt-3">
              <AlertCircle className="w-4 h-4" />
              <p>Changes to roles will affect user permissions immediately</p>
            </div>
          </div>
          <div className="flex flex-col gap-2">
            <Button 
              className="w-full sm:w-auto"
              onClick={() => navigate("/user-roles/create")}
            >
              <Plus className="w-4 h-4 mr-2" />
              Add New Role
            </Button>
            <p className="text-xs text-center text-muted-foreground">Total roles: {roles.length}</p>
          </div>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="min-w-[200px]">Role Name</TableHead>
                  <TableHead className="hidden md:table-cell">Description & Access Level</TableHead>
                  <TableHead className="min-w-[120px]">Active Users</TableHead>
                  <TableHead className="text-right min-w-[140px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {roles.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan="4" className="text-center text-muted-foreground h-32">
                      <div className="flex flex-col items-center gap-2">
                        <ShieldCheck className="w-8 h-8 text-muted-foreground/50" />
                        <p>No roles found. Click "Add New Role" to create one.</p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  roles.map((role) => (
                    <TableRow key={role.id}>
                      <TableCell className="font-medium">
                        <div className="flex items-center gap-2 sm:gap-3">
                          <div className="p-1.5 sm:p-2 bg-primary/10 rounded-lg flex-shrink-0">
                            <ShieldCheck className="w-3 h-3 sm:w-4 sm:h-4 text-primary" />
                          </div>
                          <div className="flex flex-col min-w-0">
                            <span className="text-sm sm:text-base truncate">{role.name}</span>
                            <span className="text-xs text-muted-foreground">ID: {role.id}</span>
                            <div className="md:hidden mt-1">
                              <span className="text-xs text-muted-foreground line-clamp-2">{role.description}</span>
                              <div className="mt-1">
                                {role.name === "Admin" && (
                                  <Badge variant="destructive" className="text-xs">Full Access</Badge>
                                )}
                                {role.name === "Manager" && (
                                  <Badge variant="default" className="text-xs">Limited Access</Badge>
                                )}
                                {role.name === "Support" && (
                                  <Badge variant="secondary" className="text-xs">Support Access</Badge>
                                )}
                                {role.name === "Viewer" && (
                                  <Badge variant="outline" className="text-xs">Read Only</Badge>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="hidden md:table-cell">
                        <div className="flex flex-col gap-2">
                          <span className="text-sm">{role.description}</span>
                          <div className="flex gap-1">
                            {role.name === "Admin" && (
                              <Badge variant="destructive" className="text-xs">Full Access</Badge>
                            )}
                            {role.name === "Manager" && (
                              <Badge variant="default" className="text-xs">Limited Access</Badge>
                            )}
                            {role.name === "Support" && (
                              <Badge variant="secondary" className="text-xs">Support Access</Badge>
                            )}
                            {role.name === "Viewer" && (
                              <Badge variant="outline" className="text-xs">Read Only</Badge>
                            )}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <div className="p-1 sm:p-1.5 bg-muted rounded-md flex-shrink-0">
                            <Users className="w-3 h-3 sm:w-4 sm:h-4 text-muted-foreground" />
                          </div>
                          <div className="flex flex-col min-w-0">
                            <span className="text-xs sm:text-sm font-medium">{role.users}</span>
                            <span className="text-xs text-muted-foreground">
                              {role.users === 0 ? "No users" : role.users === 1 ? "1 user" : `${role.users} users`}
                            </span>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-1">
                          <Button 
                            variant="outline" 
                            size="sm" 
                            className="h-8 px-2 text-xs"
                            onClick={() => handleEditRole(role.id)}
                          >
                            <Edit className="w-3 h-3 sm:w-3.5 sm:h-3.5 sm:mr-1" />
                            <span className="hidden sm:inline">Edit</span>
                          </Button>
                          <Button 
                            variant="ghost" 
                            size="sm"
                            className="h-8 px-2 text-xs hover:bg-destructive/10 hover:text-destructive"
                            onClick={() => handleDeleteRole(role.id)}
                          >
                            <Trash2 className="w-3 h-3 sm:w-3.5 sm:h-3.5 sm:mr-1" />
                            <span className="hidden sm:inline">Delete</span>
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}