package templates

import "fmt"

// ResetPasswordEmail generates HTML template for reset password email.
func ResetPasswordEmail(token, resetURL, appName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            background-color: #f4f4f4;
            border-radius: 8px;
            padding: 30px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .header h1 {
            color: #2c3e50;
            margin: 0;
        }
        .content {
            background-color: #ffffff;
            padding: 25px;
            border-radius: 5px;
        }
        .button {
            display: inline-block;
            background-color: #3498db;
            color: #ffffff;
            text-decoration: none;
            padding: 12px 30px;
            border-radius: 5px;
            margin: 20px 0;
            font-weight: bold;
        }
        .button:hover {
            background-color: #2980b9;
        }
        .token {
            background-color: #ecf0f1;
            padding: 10px 15px;
            border-radius: 4px;
            font-family: monospace;
            font-size: 14px;
            word-break: break-all;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            color: #7f8c8d;
            font-size: 12px;
        }
        .warning {
            background-color: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin-top: 20px;
            font-size: 13px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>%s</h1>
        </div>
        <div class="content">
            <p>Halo,</p>
            <p>Kami menerima permintaan untuk mereset password akun Anda. Klik tombol di bawah ini untuk mereset password:</p>
            
            <p style="text-align: center;">
                <a href="%s" class="button">Reset Password</a>
            </p>
            
            <p>Atau salin dan tempel link berikut ke browser Anda:</p>
            <p class="token">%s</p>
            
            <p>Atau gunakan token berikut:</p>
            <p class="token">%s</p>
            
            <div class="warning">
                <strong>⚠️ Penting:</strong> Link ini hanya berlaku selama 1 jam dan hanya dapat digunakan sekali. 
                Jika Anda tidak meminta reset password, abaikan email ini atau hubungi support jika Anda memiliki kekhawatiran.
            </div>
        </div>
        <div class="footer">
            <p>&copy; %s. All rights reserved.</p>
            <p>Email ini dikirim secara otomatis, mohon tidak membalas email ini.</p>
        </div>
    </div>
</body>
</html>
`, appName, resetURL, resetURL, token, appName)
}
