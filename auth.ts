import { IRequest, error } from "itty-router"
import { verifyTOTP } from "totp-basic"
import { Env } from "."

export async function withTOTPQueryAuth(request: IRequest, env: Env) {
  const otp = request.query['otp']
  if (typeof otp !== 'string') {
    return error(401, 'Query `otp` is needed.')
  }
  if (!await verifyTOTP(env.TOTP_SECRET, otp)) {
    return error(403, 'Invalid OTP.')
  }
}

export async function withTOTPHeaderAuth(request: IRequest, env: Env) {
  const auth = request.headers.get('Authorization')
  const [scheme, otp] = auth?.split(' ') ?? []
  if (scheme !== 'TOTP') {
    return error(401, 'Header `Authorization` with scheme `TOTP` is needed.')
  }
  if (!await verifyTOTP(env.TOTP_SECRET, otp)) {
    return error(403, 'Invalid OTP.')
  }
}