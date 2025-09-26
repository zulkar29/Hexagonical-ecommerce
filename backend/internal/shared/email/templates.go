package email

// getVerificationEmailTemplate returns the email verification template
func (s *Service) getVerificationEmailTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Account</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #ffffff; padding: 30px; border: 1px solid #dee2e6; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 0 0 8px 8px; font-size: 12px; color: #6c757d; }
        .button { display: inline-block; padding: 12px 30px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .button:hover { background-color: #0056b3; }
        .alert { background-color: #fff3cd; border: 1px solid #ffeaa7; color: #856404; padding: 15px; border-radius: 4px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{company_name}}</h1>
        </div>
        <div class="content">
            <h2>Verify Your Email Address</h2>
            <p>Thank you for signing up! Please click the button below to verify your email address and activate your account.</p>

            <div style="text-align: center;">
                <a href="{{verify_url}}" class="button">Verify Email Address</a>
            </div>

            <p>If the button doesn't work, you can also copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #007bff;">{{verify_url}}</p>

            <div class="alert">
                <strong>Note:</strong> This verification link will expire in 24 hours for security reasons.
            </div>

            <p>If you didn't create an account with us, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>&copy; {{company_name}}. All rights reserved.</p>
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`
}

// getPasswordResetEmailTemplate returns the password reset template
func (s *Service) getPasswordResetEmailTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #ffffff; padding: 30px; border: 1px solid #dee2e6; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 0 0 8px 8px; font-size: 12px; color: #6c757d; }
        .button { display: inline-block; padding: 12px 30px; background-color: #dc3545; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .button:hover { background-color: #c82333; }
        .alert { background-color: #f8d7da; border: 1px solid #f5c6cb; color: #721c24; padding: 15px; border-radius: 4px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{company_name}}</h1>
        </div>
        <div class="content">
            <h2>Reset Your Password</h2>
            <p>We received a request to reset your password. Click the button below to create a new password.</p>

            <div style="text-align: center;">
                <a href="{{reset_url}}" class="button">Reset Password</a>
            </div>

            <p>If the button doesn't work, you can also copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #dc3545;">{{reset_url}}</p>

            <div class="alert">
                <strong>Important:</strong> This password reset link will expire in 1 hour for security reasons.
            </div>

            <p>If you didn't request a password reset, please ignore this email or contact support if you have concerns.</p>
        </div>
        <div class="footer">
            <p>&copy; {{company_name}}. All rights reserved.</p>
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`
}

// getOrderConfirmationEmailTemplate returns the order confirmation template
func (s *Service) getOrderConfirmationEmailTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Order Confirmation</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #28a745; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #ffffff; padding: 30px; border: 1px solid #dee2e6; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 0 0 8px 8px; font-size: 12px; color: #6c757d; }
        .order-info { background-color: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .success { color: #28a745; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{company_name}}</h1>
            <p>Order Confirmation</p>
        </div>
        <div class="content">
            <h2 class="success">✓ Order Confirmed!</h2>
            <p>Thank you for your order! We've received your order and are processing it now.</p>

            <div class="order-info">
                <h3>Order Details</h3>
                <p><strong>Order Number:</strong> #{{order_number}}</p>
                <p><strong>Order Date:</strong> {{order_date}}</p>
                <p><strong>Status:</strong> Processing</p>
            </div>

            <h3>What's Next?</h3>
            <ul>
                <li>You'll receive an email when your order ships</li>
                <li>Track your order status in your account dashboard</li>
                <li>Contact support if you have any questions</li>
            </ul>

            <p>We appreciate your business and will send you a shipping confirmation email with tracking information once your order is on its way.</p>
        </div>
        <div class="footer">
            <p>&copy; {{company_name}}. All rights reserved.</p>
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`
}

// getPaymentSuccessEmailTemplate returns the payment success template
func (s *Service) getPaymentSuccessEmailTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Payment Successful</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #28a745; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #ffffff; padding: 30px; border: 1px solid #dee2e6; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 0 0 8px 8px; font-size: 12px; color: #6c757d; }
        .payment-info { background-color: #d4edda; border: 1px solid #c3e6cb; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .success { color: #28a745; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{company_name}}</h1>
            <p>Payment Confirmation</p>
        </div>
        <div class="content">
            <h2 class="success">✓ Payment Successful!</h2>
            <p>Great news! Your payment has been successfully processed.</p>

            <div class="payment-info">
                <h3>Payment Details</h3>
                <p><strong>Order Number:</strong> #{{order_number}}</p>
                <p><strong>Payment Status:</strong> Completed</p>
                <p><strong>Payment Date:</strong> {{payment_date}}</p>
            </div>

            <p>Your order is now being prepared for shipment. You'll receive a shipping confirmation email with tracking information once your order is dispatched.</p>

            <p>Thank you for choosing {{company_name}}!</p>
        </div>
        <div class="footer">
            <p>&copy; {{company_name}}. All rights reserved.</p>
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`
}

// getPaymentFailedEmailTemplate returns the payment failed template
func (s *Service) getPaymentFailedEmailTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Payment Failed</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #dc3545; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #ffffff; padding: 30px; border: 1px solid #dee2e6; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 0 0 8px 8px; font-size: 12px; color: #6c757d; }
        .payment-info { background-color: #f8d7da; border: 1px solid #f5c6cb; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .error { color: #dc3545; font-weight: bold; }
        .button { display: inline-block; padding: 12px 30px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{company_name}}</h1>
            <p>Payment Issue</p>
        </div>
        <div class="content">
            <h2 class="error">Payment Failed</h2>
            <p>We encountered an issue processing your payment for order #{{order_number}}.</p>

            <div class="payment-info">
                <h3>Issue Details</h3>
                <p><strong>Order Number:</strong> #{{order_number}}</p>
                <p><strong>Reason:</strong> {{reason}}</p>
            </div>

            <h3>What You Can Do:</h3>
            <ul>
                <li>Check your payment method details</li>
                <li>Ensure sufficient funds are available</li>
                <li>Try a different payment method</li>
                <li>Contact your bank if the issue persists</li>
            </ul>

            <div style="text-align: center;">
                <a href="{{retry_url}}" class="button">Retry Payment</a>
            </div>

            <p>Your order is currently on hold. Please complete the payment within 24 hours to avoid cancellation.</p>

            <p>If you continue to experience issues, please contact our support team.</p>
        </div>
        <div class="footer">
            <p>&copy; {{company_name}}. All rights reserved.</p>
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`
}