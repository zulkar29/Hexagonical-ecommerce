import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { useNavigate, useParams } from "react-router-dom";
import { 
  ArrowLeft, 
  Plus,
  Minus,
  Receipt,
  Upload,
  Calendar,
  DollarSign,
  FileText,
  User,
  AlertCircle,
  Loader2
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "@/components/ui/alert";
// import { pettyCashApi } from "@/lib/api";

const EditPettyCashTransaction = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingData, setIsLoadingData] = useState(true);
  const [transactionType, setTransactionType] = useState("out");
  const [receiptFile, setReceiptFile] = useState(null);
  const [existingReceipt, setExistingReceipt] = useState(null);

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    reset,
    formState: { errors }
  } = useForm({
    defaultValues: {
      type: "out",
      amount: "",
      category: "",
      description: "",
      date: "",
      handledBy: "",
      requiresApproval: true,
      notes: ""
    }
  });

  const categories = [
    "Office Supplies",
    "Travel",
    "Refreshments", 
    "Maintenance",
    "Utilities",
    "Marketing",
    "Training",
    "Emergency",
    "Miscellaneous",
    "Replenishment"
  ];

  const amount = watch("amount");
  const selectedCategory = watch("category");

  // Load transaction data
  useEffect(() => {
    const loadTransaction = async () => {
      try {
        // TODO: Replace with actual API call
        // const response = await pettyCashApi.getTransaction(id);
        
        // Mock data
        const mockTransaction = {
          id: parseInt(id),
          type: "out",
          amount: 89.50,
          category: "Office Supplies",
          description: "Printer paper and ink cartridges",
          date: "2024-03-18",
          handledBy: "John Admin",
          requiresApproval: false,
          notes: "Urgent purchase for marketing presentation",
          approved: true,
          receipt: "receipt-001234.pdf"
        };

        // Populate form with existing data
        reset({
          type: mockTransaction.type,
          amount: mockTransaction.amount.toString(),
          category: mockTransaction.category,
          description: mockTransaction.description,
          date: mockTransaction.date,
          handledBy: mockTransaction.handledBy,
          requiresApproval: !mockTransaction.approved,
          notes: mockTransaction.notes || ""
        });

        setTransactionType(mockTransaction.type);
        setExistingReceipt(mockTransaction.receipt);
        
      } catch (error) {
        toast.error("Failed to load transaction details");
        console.error(error);
        navigate("/petty-cash");
      } finally {
        setIsLoadingData(false);
      }
    };

    loadTransaction();
  }, [id, reset, navigate]);

  const onSubmit = async (data) => {
    try {
      setIsLoading(true);
      
      // TODO: Implement API call to update transaction
      const transactionData = {
        ...data,
        type: transactionType,
        amount: parseFloat(data.amount),
        receipt: receiptFile ? receiptFile.name : existingReceipt,
        approved: !data.requiresApproval
      };

      console.log("Updating transaction:", transactionData);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      toast.success("Transaction updated successfully!");
      navigate("/petty-cash");
    } catch (error) {
      toast.error("Failed to update transaction. Please try again.");
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      if (file.size > 5 * 1024 * 1024) { // 5MB limit
        toast.error("File size must be less than 5MB");
        return;
      }
      setReceiptFile(file);
      toast.success("Receipt uploaded successfully");
    }
  };

  const getCurrentBalance = () => {
    // Mock current balance - in real app, fetch from API
    return 2450.75;
  };

  const getEstimatedBalance = () => {
    const currentBalance = getCurrentBalance();
    const transactionAmount = parseFloat(amount) || 0;
    
    if (transactionType === "out") {
      return currentBalance - transactionAmount;
    } else {
      return currentBalance + transactionAmount;
    }
  };

  if (isLoadingData) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto" />
          <p className="mt-2 text-muted-foreground">Loading transaction details...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6 max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/petty-cash")}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Edit Petty Cash Transaction</h1>
          <p className="text-muted-foreground">
            Update transaction details and information
          </p>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Form */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Transaction Details</CardTitle>
              <CardDescription>
                Update the petty cash transaction information
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                {/* Transaction Type */}
                <div className="space-y-3">
                  <Label>Transaction Type</Label>
                  <div className="flex gap-4">
                    <Button
                      type="button"
                      variant={transactionType === "out" ? "default" : "outline"}
                      onClick={() => {
                        setTransactionType("out");
                        setValue("type", "out");
                      }}
                      className="flex-1"
                    >
                      <Minus className="mr-2 h-4 w-4" />
                      Cash Out (Expense)
                    </Button>
                    <Button
                      type="button"
                      variant={transactionType === "in" ? "default" : "outline"}
                      onClick={() => {
                        setTransactionType("in");
                        setValue("type", "in");
                      }}
                      className="flex-1"
                    >
                      <Plus className="mr-2 h-4 w-4" />
                      Cash In (Replenishment)
                    </Button>
                  </div>
                </div>

                <Separator />

                {/* Amount and Date */}
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="amount">Amount *</Label>
                    <div className="relative">
                      <DollarSign className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                      <Input
                        id="amount"
                        type="number"
                        step="0.01"
                        placeholder="0.00"
                        className={`pl-10 ${errors.amount ? "border-red-500" : ""}`}
                        {...register("amount", { 
                          required: "Amount is required",
                          min: { value: 0.01, message: "Amount must be greater than 0" }
                        })}
                      />
                    </div>
                    {errors.amount && (
                      <p className="text-sm text-red-500">{errors.amount.message}</p>
                    )}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="date">Date *</Label>
                    <div className="relative">
                      <Calendar className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                      <Input
                        id="date"
                        type="date"
                        className={`pl-10 ${errors.date ? "border-red-500" : ""}`}
                        {...register("date", { required: "Date is required" })}
                      />
                    </div>
                    {errors.date && (
                      <p className="text-sm text-red-500">{errors.date.message}</p>
                    )}
                  </div>
                </div>

                {/* Category */}
                <div className="space-y-2">
                  <Label htmlFor="category">Category *</Label>
                  <Select 
                    value={watch("category")}
                    onValueChange={(value) => setValue("category", value)}
                  >
                    <SelectTrigger className={errors.category ? "border-red-500" : ""}>
                      <SelectValue placeholder="Select a category" />
                    </SelectTrigger>
                    <SelectContent>
                      {categories.map((category) => (
                        <SelectItem key={category} value={category}>
                          {category}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {errors.category && (
                    <p className="text-sm text-red-500">{errors.category.message}</p>
                  )}
                </div>

                {/* Description */}
                <div className="space-y-2">
                  <Label htmlFor="description">Description *</Label>
                  <div className="relative">
                    <FileText className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                    <Textarea
                      id="description"
                      placeholder="Describe the transaction..."
                      className={`pl-10 ${errors.description ? "border-red-500" : ""}`}
                      rows={3}
                      {...register("description", { 
                        required: "Description is required",
                        minLength: { value: 10, message: "Description must be at least 10 characters" }
                      })}
                    />
                  </div>
                  {errors.description && (
                    <p className="text-sm text-red-500">{errors.description.message}</p>
                  )}
                </div>

                {/* Handled By */}
                <div className="space-y-2">
                  <Label htmlFor="handledBy">Handled By *</Label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                      id="handledBy"
                      placeholder="Person responsible for this transaction"
                      className={`pl-10 ${errors.handledBy ? "border-red-500" : ""}`}
                      {...register("handledBy", { required: "Handler name is required" })}
                    />
                  </div>
                  {errors.handledBy && (
                    <p className="text-sm text-red-500">{errors.handledBy.message}</p>
                  )}
                </div>

                {/* Receipt Upload */}
                <div className="space-y-2">
                  <Label htmlFor="receipt">Receipt/Document</Label>
                  <div className="border-2 border-dashed border-muted-foreground/25 rounded-lg p-6 text-center">
                    <Receipt className="mx-auto h-8 w-8 text-muted-foreground mb-2" />
                    <div className="space-y-2">
                      {existingReceipt && !receiptFile && (
                        <p className="text-sm text-green-600">
                          Current receipt: {existingReceipt}
                        </p>
                      )}
                      <p className="text-sm text-muted-foreground">
                        {receiptFile ? receiptFile.name : "Upload new receipt or supporting document"}
                      </p>
                      <input
                        type="file"
                        id="receipt"
                        accept="image/*,.pdf"
                        onChange={handleFileUpload}
                        className="hidden"
                      />
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => document.getElementById("receipt").click()}
                      >
                        <Upload className="mr-2 h-4 w-4" />
                        {receiptFile ? "Change File" : existingReceipt ? "Replace File" : "Choose File"}
                      </Button>
                    </div>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Supported formats: JPG, PNG, PDF (max 5MB)
                  </p>
                </div>

                {/* Additional Options */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="space-y-1">
                      <Label htmlFor="requiresApproval">Requires Approval</Label>
                      <p className="text-sm text-muted-foreground">
                        Transaction will be pending until approved
                      </p>
                    </div>
                    <Switch
                      id="requiresApproval"
                      checked={watch("requiresApproval")}
                      onCheckedChange={(checked) => setValue("requiresApproval", checked)}
                    />
                  </div>
                </div>

                {/* Additional Notes */}
                <div className="space-y-2">
                  <Label htmlFor="notes">Additional Notes</Label>
                  <Textarea
                    id="notes"
                    placeholder="Any additional information..."
                    rows={2}
                    {...register("notes")}
                  />
                </div>

                {/* Submit Buttons */}
                <div className="flex justify-end gap-3 pt-6">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => navigate("/petty-cash")}
                  >
                    Cancel
                  </Button>
                  <Button type="submit" disabled={isLoading}>
                    {isLoading ? "Updating..." : "Update Transaction"}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>

        {/* Summary Sidebar */}
        <div className="space-y-6">
          {/* Current Balance */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Balance Summary</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between">
                <span className="text-sm text-muted-foreground">Current Balance:</span>
                <span className="font-medium">${getCurrentBalance().toFixed(2)}</span>
              </div>
              
              {amount && (
                <>
                  <div className="flex justify-between">
                    <span className="text-sm text-muted-foreground">Transaction:</span>
                    <span className={`font-medium ${transactionType === 'in' ? 'text-green-600' : 'text-red-600'}`}>
                      {transactionType === 'in' ? '+' : '-'}${parseFloat(amount).toFixed(2)}
                    </span>
                  </div>
                  
                  <Separator />
                  
                  <div className="flex justify-between">
                    <span className="text-sm font-medium">Estimated Balance:</span>
                    <span className="font-bold">${getEstimatedBalance().toFixed(2)}</span>
                  </div>
                  
                  {transactionType === 'out' && getEstimatedBalance() < 100 && (
                    <Alert>
                      <AlertCircle className="h-4 w-4" />
                      <AlertTitle>Low Balance Warning</AlertTitle>
                      <AlertDescription>
                        This transaction will result in a low balance. Consider replenishing petty cash.
                      </AlertDescription>
                    </Alert>
                  )}
                </>
              )}
            </CardContent>
          </Card>

          {/* Category Guidelines */}
          {selectedCategory && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Category Guidelines</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2 text-sm">
                  {selectedCategory === "Office Supplies" && (
                    <p>Includes stationery, printing materials, basic office equipment under $100.</p>
                  )}
                  {selectedCategory === "Travel" && (
                    <p>Local transportation, parking fees, emergency travel expenses.</p>
                  )}
                  {selectedCategory === "Refreshments" && (
                    <p>Office coffee, water, meeting refreshments, staff treats.</p>
                  )}
                  {selectedCategory === "Maintenance" && (
                    <p>Minor repairs, cleaning supplies, basic maintenance items.</p>
                  )}
                  {selectedCategory === "Emergency" && (
                    <p>Urgent expenses that cannot wait for normal approval process.</p>
                  )}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
};

export default EditPettyCashTransaction;