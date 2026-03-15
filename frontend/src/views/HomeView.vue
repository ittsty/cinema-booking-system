<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { signOut } from "firebase/auth";
import { auth } from "../firebase/config";
import { apiGet } from "../api/client";

type Showtime = {
  id: string;
  movie_id: string;
  movie_title: string;
  start_time: string;
  theater: string;
};

const router = useRouter();
const showtimes = ref<Showtime[]>([]);
const loading = ref(false);

async function loadShowtimes() {
  try {
    loading.value = true;
    const res = await apiGet("/showtimes");
    showtimes.value = Array.isArray(res) ? res : res.data || [];
  } catch (error) {
    console.error(error);
    alert(String(error));
  } finally {
    loading.value = false;
  }
}

async function logout() {
  await signOut(auth);
  router.push("/login");
}

function goToSeatMap(showtimeId: string) {
  router.push(`/showtimes/${showtimeId}/seats`);
}

function formatDate(value: string) {
  return new Date(value).toLocaleString();
}

onMounted(loadShowtimes);
</script>

<template>
  <div class="min-h-screen bg-slate-100 p-6">
    <div class="mx-auto max-w-7xl space-y-6">
      <div
        class="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm md:flex-row md:items-center md:justify-between"
      >
        <div>
          <h1 class="text-3xl font-bold text-slate-900">Showtimes</h1>
          <p class="mt-1 text-sm text-slate-500">
            Select a showtime to view seats
          </p>
        </div>

        <div class="flex gap-3">
          <button
            @click="router.push('/admin')"
            class="rounded-xl border border-slate-300 bg-white px-4 py-2 font-medium text-slate-700 hover:bg-slate-50"
          >
            Admin
          </button>
          <button
            @click="logout"
            class="rounded-xl bg-slate-900 px-4 py-2 font-medium text-white hover:bg-slate-800"
          >
            Logout
          </button>
        </div>
      </div>

      <div
        v-if="loading"
        class="rounded-2xl bg-white p-10 text-center text-slate-500 shadow-sm"
      >
        Loading...
      </div>

      <div
        v-else-if="showtimes.length === 0"
        class="rounded-2xl bg-white p-10 text-center text-slate-500 shadow-sm"
      >
        No showtimes found
      </div>

      <div v-else class="grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
        <div
          v-for="showtime in showtimes"
          :key="showtime.id"
          class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm transition hover:-translate-y-1 hover:shadow-md"
        >
          <div class="space-y-2">
            <h2 class="text-xl font-semibold text-slate-900">
              {{ showtime.movie_title }}
            </h2>

            <p class="text-sm text-slate-500">
              Showtime ID: {{ showtime.id }}
            </p>

            <p class="text-sm text-slate-600">
              Theater: {{ showtime.theater || "-" }}
            </p>

            <p class="text-sm text-slate-600">
              Time: {{ formatDate(showtime.start_time) }}
            </p>
          </div>

          <button
            @click="goToSeatMap(showtime.id)"
            class="mt-5 w-full rounded-xl bg-slate-900 px-4 py-3 font-medium text-white hover:bg-slate-800"
          >
            View Seats
          </button>
        </div>
      </div>
    </div>
  </div>
</template>