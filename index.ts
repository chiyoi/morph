import { Router, error, json } from "itty-router"
import { deleteFile, getFile, listFiles, putFile, withBucket } from "./r2"
import { withTOTPHeaderAuth, withTOTPQueryAuth } from "./auth"
import { deleteKey, getKey, listKeys, putKey, withNamespace } from "./kv"

export default {
  fetch: (request: Request, env: Env, ctx: ExecutionContext) => router()
    .handle(request, env, ctx)
    .catch(error)
}

function router() {
  const router = Router()
  router.all('/ping', () => json('Pong!'))
  router.all('/version', (_, env: Env) => json(env.VERSION))
  router.all('/', (_, env: Env) => json(env.VERSION))

  router.get('/files', withTOTPQueryAuth, withBucket('files'), listFiles)
  router.all('/files', () => error(405, 'Readonly endpoint (List Files).'))
  router.get('/files/:filename', withTOTPQueryAuth, withBucket('files'), getFile)
  router.put('/files/:filename', withTOTPQueryAuth, withBucket('files'), putFile)
  router.delete('/files/:filename', withTOTPQueryAuth, withBucket('files'), deleteFile)
  router.all('/files/:filename', () => error(405, 'Method not allowed.'))

  router.get('/assets/:filename', withBucket('assets'), getFile)
  router.all('/assets/:filename', () => error(405, 'Assets is readonly.'))
  router.all('/assets', () => error(403, 'Cannot list assets.'))

  router.get('/table/keys', withTOTPHeaderAuth, withNamespace('table'), listKeys)
  router.all('/table/keys', () => error(405, 'Readonly endpoint (List Table Keys).'))
  router.get('/table/:key', withTOTPHeaderAuth, withNamespace('table'), getKey)
  router.put('/table/:key', withTOTPHeaderAuth, withNamespace('table'), putKey)
  router.delete('/table/:key', withTOTPHeaderAuth, withNamespace('table'), deleteKey)
  router.all('/table/:key', () => error(405, 'Method not allowed.'))
  router.all('/table', () => error(400, 'Sub-path is needed.'))

  router.all('*', () => error(404, 'Invalid path.'))
  return router
}


export interface Env {
  VERSION: string,
  TOTP_SECRET: string,

  table: KVNamespace,

  files: R2Bucket,
  assets: R2Bucket,
}