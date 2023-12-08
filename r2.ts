import { IRequest, error, json } from "itty-router"
import { Env } from "."

export function withBucket(bucket: 'files' | 'assets') {
  return (request: IRequest & { bucket: R2Bucket }, env: Env) => {
    request.bucket = env[bucket]
  }
}

export async function listFiles(request: IRequest & { bucket: R2Bucket }) {
  const { bucket } = request
  const files = await bucket.list()
  return json(files.objects.map(file => file.key))
}

export async function getFile(request: IRequest & { bucket: R2Bucket }) {
  const { params: { filename }, bucket } = request
  const file = await bucket.get(filename)
  if (file === null) {
    return error(404, 'No such file.')
  }
  const headers = new Headers()
  file.writeHttpMetadata(headers)
  return new Response(file.body, { headers })
}

export async function putFile(request: IRequest & { bucket: R2Bucket, params: { filename: string } }) {
  const { bucket, params: { filename } } = request
  await bucket.put(filename, request.body)
  return json(`Put ${filename}.`)
}

export async function deleteFile(request: IRequest & { bucket: R2Bucket, params: { filename: string } }) {
  const { bucket, params: { filename } } = request
  const head = await bucket.head(filename)
  if (head === null) {
    return error(404, 'No such file.')
  }
  await bucket.delete(filename)
  return json(`Deleted ${filename}.`)
}
