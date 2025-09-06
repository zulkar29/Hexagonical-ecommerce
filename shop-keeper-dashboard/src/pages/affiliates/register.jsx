import { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { 
  ArrowLeft, 
  User, 
  Mail, 
  Phone, 
  Globe, 
  FileText, 
  DollarSign,
  Users,
  TrendingUp,
  Shield,
  CheckCircle
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
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Checkbox } from "@/components/ui/checkbox";

const AffiliateRegister = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [currentStep, setCurrentStep] = useState(1);

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors }
  } = useForm({
    defaultValues: {
      personal: {
        firstName: "",
        lastName: "",
        email: "",
        phone: "",
        country: "",
        state: "",
        city: ""
      },
      business: {
        businessName: "",
        website: "",
        socialMedia: "",
        audienceSize: "",
        niche: "",
        experience: "",
        promotionMethods: []
      },
      agreement: {
        termsAccepted: false,
        dataProcessingAccepted: false
      }
    }
  });

  const onSubmit = async (data) => {
    if (currentStep < 3) {
      setCurrentStep(currentStep + 1);
      return;
    }

    try {
      setIsLoading(true);
      // TODO: Implement affiliate registration API call
      console.log("Registration data:", data);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      toast.success("Application submitted successfully! We'll review it within 24 hours.");
      navigate("/affiliates");
    } catch (error) {
      toast.error("Failed to submit application. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  const benefits = [
    {
      icon: <DollarSign className="h-8 w-8 text-green-500" />,
      title: "Competitive Commissions",
      description: "Earn up to 15% commission on every sale you refer"
    },
    {
      icon: <Users className="h-8 w-8 text-blue-500" />,
      title: "Marketing Support",
      description: "Access to marketing materials, banners, and promotional content"
    },
    {
      icon: <TrendingUp className="h-8 w-8 text-purple-500" />,
      title: "Real-time Analytics",
      description: "Track your performance with detailed analytics and reporting"
    },
    {
      icon: <Shield className="h-8 w-8 text-orange-500" />,
      title: "Reliable Payouts",
      description: "Monthly payments via PayPal, bank transfer, or other methods"
    }
  ];

  const stepTitles = [
    "Personal Information",
    "Business Details", 
    "Terms & Agreement"
  ];

  const renderStepIndicator = () => (
    <div className="flex items-center justify-center mb-8">
      {[1, 2, 3].map((step, index) => (
        <div key={step} className="flex items-center">
          <div className={`
            flex items-center justify-center w-8 h-8 rounded-full border-2 text-sm font-medium
            ${currentStep >= step 
              ? 'bg-primary text-primary-foreground border-primary' 
              : 'bg-background text-muted-foreground border-muted-foreground'
            }
          `}>
            {currentStep > step ? <CheckCircle className="w-5 h-5" /> : step}
          </div>
          {index < 2 && (
            <div className={`w-16 h-0.5 mx-2 ${
              currentStep > step + 1 ? 'bg-primary' : 'bg-muted-foreground'
            }`} />
          )}
        </div>
      ))}
    </div>
  );

  return (
    <div className="min-h-screen bg-background p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/affiliates")}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold">Join Our Affiliate Program</h1>
            <p className="text-muted-foreground">
              Partner with us and start earning commissions today
            </p>
          </div>
        </div>

        {/* Benefits Section - Only show on first step */}
        {currentStep === 1 && (
          <Card className="mb-8">
            <CardHeader>
              <CardTitle className="text-center">Why Join Our Affiliate Program?</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
                {benefits.map((benefit, index) => (
                  <div key={index} className="text-center space-y-3">
                    <div className="flex justify-center">{benefit.icon}</div>
                    <h3 className="font-semibold">{benefit.title}</h3>
                    <p className="text-sm text-muted-foreground">{benefit.description}</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}

        {/* Registration Form */}
        <Card>
          <CardHeader>
            <div className="text-center">
              <CardTitle>Application Form</CardTitle>
              <CardDescription>
                Step {currentStep} of 3: {stepTitles[currentStep - 1]}
              </CardDescription>
            </div>
          </CardHeader>
          <CardContent>
            {renderStepIndicator()}

            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              {/* Step 1: Personal Information */}
              {currentStep === 1 && (
                <div className="space-y-6">
                  <div className="grid md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="firstName">First Name *</Label>
                      <Input
                        id="firstName"
                        {...register("personal.firstName", { 
                          required: "First name is required" 
                        })}
                        className={errors.personal?.firstName ? "border-red-500" : ""}
                      />
                      {errors.personal?.firstName && (
                        <p className="text-sm text-red-500">{errors.personal.firstName.message}</p>
                      )}
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="lastName">Last Name *</Label>
                      <Input
                        id="lastName"
                        {...register("personal.lastName", { 
                          required: "Last name is required" 
                        })}
                        className={errors.personal?.lastName ? "border-red-500" : ""}
                      />
                      {errors.personal?.lastName && (
                        <p className="text-sm text-red-500">{errors.personal.lastName.message}</p>
                      )}
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="email">Email Address *</Label>
                    <Input
                      id="email"
                      type="email"
                      {...register("personal.email", { 
                        required: "Email is required",
                        pattern: {
                          value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                          message: "Invalid email address"
                        }
                      })}
                      className={errors.personal?.email ? "border-red-500" : ""}
                    />
                    {errors.personal?.email && (
                      <p className="text-sm text-red-500">{errors.personal.email.message}</p>
                    )}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="phone">Phone Number</Label>
                    <Input
                      id="phone"
                      {...register("personal.phone")}
                      placeholder="Optional"
                    />
                  </div>

                  <div className="grid md:grid-cols-3 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="country">Country *</Label>
                      <Select onValueChange={(value) => setValue("personal.country", value)}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select country" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="us">United States</SelectItem>
                          <SelectItem value="ca">Canada</SelectItem>
                          <SelectItem value="uk">United Kingdom</SelectItem>
                          <SelectItem value="au">Australia</SelectItem>
                          <SelectItem value="de">Germany</SelectItem>
                          <SelectItem value="fr">France</SelectItem>
                          <SelectItem value="other">Other</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="state">State/Province</Label>
                      <Input
                        id="state"
                        {...register("personal.state")}
                        placeholder="Optional"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="city">City</Label>
                      <Input
                        id="city"
                        {...register("personal.city")}
                        placeholder="Optional"
                      />
                    </div>
                  </div>
                </div>
              )}

              {/* Step 2: Business Details */}
              {currentStep === 2 && (
                <div className="space-y-6">
                  <div className="space-y-2">
                    <Label htmlFor="businessName">Business/Brand Name</Label>
                    <Input
                      id="businessName"
                      {...register("business.businessName")}
                      placeholder="Your business or personal brand name"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="website">Website URL *</Label>
                    <Input
                      id="website"
                      type="url"
                      {...register("business.website", { 
                        required: "Website URL is required" 
                      })}
                      placeholder="https://your-website.com"
                      className={errors.business?.website ? "border-red-500" : ""}
                    />
                    {errors.business?.website && (
                      <p className="text-sm text-red-500">{errors.business.website.message}</p>
                    )}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="socialMedia">Primary Social Media Profile</Label>
                    <Input
                      id="socialMedia"
                      {...register("business.socialMedia")}
                      placeholder="Instagram, YouTube, TikTok, etc."
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="audienceSize">Audience Size *</Label>
                    <Select onValueChange={(value) => setValue("business.audienceSize", value)}>
                      <SelectTrigger>
                        <SelectValue placeholder="Select audience size" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="0-1k">0 - 1,000</SelectItem>
                        <SelectItem value="1k-10k">1,000 - 10,000</SelectItem>
                        <SelectItem value="10k-50k">10,000 - 50,000</SelectItem>
                        <SelectItem value="50k-100k">50,000 - 100,000</SelectItem>
                        <SelectItem value="100k+">100,000+</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="niche">Your Niche/Industry *</Label>
                    <Select onValueChange={(value) => setValue("business.niche", value)}>
                      <SelectTrigger>
                        <SelectValue placeholder="Select your niche" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="fashion">Fashion & Style</SelectItem>
                        <SelectItem value="tech">Technology</SelectItem>
                        <SelectItem value="fitness">Health & Fitness</SelectItem>
                        <SelectItem value="lifestyle">Lifestyle</SelectItem>
                        <SelectItem value="business">Business</SelectItem>
                        <SelectItem value="food">Food & Cooking</SelectItem>
                        <SelectItem value="travel">Travel</SelectItem>
                        <SelectItem value="other">Other</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="experience">Affiliate Marketing Experience</Label>
                    <Textarea
                      id="experience"
                      {...register("business.experience")}
                      placeholder="Tell us about your experience with affiliate marketing..."
                      rows={4}
                    />
                  </div>
                </div>
              )}

              {/* Step 3: Terms & Agreement */}
              {currentStep === 3 && (
                <div className="space-y-6">
                  <div className="text-center space-y-4">
                    <h3 className="text-xl font-semibold">Almost Done!</h3>
                    <p className="text-muted-foreground">
                      Please review and accept our terms to complete your application.
                    </p>
                  </div>

                  <Card className="border-dashed">
                    <CardContent className="pt-6">
                      <h4 className="font-semibold mb-3">Commission Structure</h4>
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <span>Sales $0 - $1,000:</span>
                          <Badge variant="secondary">8% Commission</Badge>
                        </div>
                        <div className="flex justify-between">
                          <span>Sales $1,001 - $5,000:</span>
                          <Badge variant="secondary">10% Commission</Badge>
                        </div>
                        <div className="flex justify-between">
                          <span>Sales $5,001+:</span>
                          <Badge variant="secondary">15% Commission</Badge>
                        </div>
                      </div>
                      <Separator className="my-4" />
                      <p className="text-xs text-muted-foreground">
                        Payouts are processed monthly on the 1st of each month for the previous month's commissions.
                        Minimum payout threshold is $50.
                      </p>
                    </CardContent>
                  </Card>

                  <div className="space-y-4">
                    <div className="flex items-start space-x-3">
                      <Checkbox
                        id="terms"
                        {...register("agreement.termsAccepted", {
                          required: "You must accept the terms and conditions"
                        })}
                      />
                      <div className="space-y-1">
                        <Label htmlFor="terms" className="text-sm leading-relaxed">
                          I agree to the <span className="text-primary underline cursor-pointer">Affiliate Terms and Conditions</span> and <span className="text-primary underline cursor-pointer">Privacy Policy</span>
                        </Label>
                        {errors.agreement?.termsAccepted && (
                          <p className="text-sm text-red-500">{errors.agreement.termsAccepted.message}</p>
                        )}
                      </div>
                    </div>

                    <div className="flex items-start space-x-3">
                      <Checkbox
                        id="dataProcessing"
                        {...register("agreement.dataProcessingAccepted", {
                          required: "You must consent to data processing"
                        })}
                      />
                      <div className="space-y-1">
                        <Label htmlFor="dataProcessing" className="text-sm leading-relaxed">
                          I consent to the processing of my personal data for affiliate program management and communication purposes
                        </Label>
                        {errors.agreement?.dataProcessingAccepted && (
                          <p className="text-sm text-red-500">{errors.agreement.dataProcessingAccepted.message}</p>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Navigation Buttons */}
              <div className="flex justify-between pt-6">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => currentStep > 1 ? setCurrentStep(currentStep - 1) : navigate("/affiliates")}
                >
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  {currentStep > 1 ? "Previous" : "Back"}
                </Button>
                
                <Button
                  type="submit"
                  disabled={isLoading}
                  className="min-w-[120px]"
                >
                  {isLoading ? (
                    "Submitting..."
                  ) : currentStep < 3 ? (
                    "Next Step"
                  ) : (
                    "Submit Application"
                  )}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default AffiliateRegister;