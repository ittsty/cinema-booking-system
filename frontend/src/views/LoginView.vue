<script setup lang="ts">
import { ref } from "vue";
import {
  signInWithPopup,
  signInWithEmailAndPassword,
} from "firebase/auth";
import { useRouter } from "vue-router";
import { auth, googleProvider } from "../firebase/config";

const router = useRouter();

const email = ref("");
const password = ref("");
const loading = ref(false);

async function loginWithGoogle() {
  try {
    loading.value = true;
    await signInWithPopup(auth, googleProvider);
    router.push("/");
  } catch (error) {
    console.error(error);
    alert("Google login failed");
  } finally {
    loading.value = false;
  }
}

async function loginWithEmail() {
  try {
    loading.value = true;
    await signInWithEmailAndPassword(auth, email.value, password.value);
    router.push("/");
  } catch (error) {
    console.error(error);
    alert("Email login failed");
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-100 flex items-center justify-center p-6">
    <div class="w-full max-w-md rounded-2xl bg-white shadow-lg border border-slate-200 p-8">
      <div class="mb-6 text-center">
        <h1 class="text-3xl font-bold text-slate-900">Cinema Booking</h1>
        <p class="mt-2 text-sm text-slate-500">
          Admin ใช้ email/password, user ใช้ Google login
        </p>
      </div>

      <div class="space-y-4">
        <input
          v-model="email"
          type="email"
          placeholder="Admin email"
          class="w-full rounded-xl border border-slate-300 px-4 py-3 outline-none focus:border-slate-500"
        />

        <input
          v-model="password"
          type="password"
          placeholder="Password"
          class="w-full rounded-xl border border-slate-300 px-4 py-3 outline-none focus:border-slate-500"
        />

        <button
          @click="loginWithEmail"
          :disabled="loading"
          class="w-full rounded-xl bg-slate-900 px-4 py-3 font-medium text-white hover:bg-slate-800 disabled:opacity-50"
        >
          Login with Email
        </button>

        <div class="text-center text-sm text-slate-400">or</div>

        <button
          @click="loginWithGoogle"
          :disabled="loading"
          class="w-full rounded-xl border border-slate-300 bg-white px-4 py-3 font-medium text-slate-700 hover:bg-slate-50 disabled:opacity-50"
        >
          Login with Google
        </button>
      </div>
    </div>
  </div>
</template>