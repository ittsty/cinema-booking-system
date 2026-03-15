import { auth } from "../firebase/config";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

async function authHeaders() {
  const user = auth.currentUser;
  const token = user ? await user.getIdToken() : null;

  return {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
}

export async function apiGet(path: string) {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    headers: await authHeaders(),
  });

  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function apiPost(path: string, body?: unknown) {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: await authHeaders(),
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!res.ok) throw new Error(await res.text());
  return res.json();
}