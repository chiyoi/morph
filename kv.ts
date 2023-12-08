import { IRequest, json } from "itty-router"
import { Env } from "."

export function withNamespace(namespace: 'table') {
  return (request: IRequest & { namespace: KVNamespace }, env: Env) => {
    request.namespace = env[namespace]
  }
}

export async function listKeys(request: IRequest & { namespace: KVNamespace }) {
  const { namespace } = request
  const list = await namespace.list()
  return json(list.keys.map(key => key.name))
}

export async function getKey(request: IRequest & { namespace: KVNamespace, params: { key: string } }) {
  const { namespace, params: { key } } = request
  return json(await namespace.get(key))
}

export async function putKey(request: IRequest & { namespace: KVNamespace, params: { key: string } }) {
  const { namespace, params: { key } } = request
  await namespace.put(key, request.body ?? '')
  return json(`Put ${key}.`)
}

export async function deleteKey(request: IRequest & { namespace: KVNamespace, params: { key: string } }) {
  const { namespace, params: { key } } = request
  await namespace.delete(key)
  return json(`Deleted ${key}.`)
}
