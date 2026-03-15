<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { apiGet } from "../api/client";

const router = useRouter();

const bookings = ref<any[]>([]);
const logs = ref<any[]>([]);
const loading = ref(false);

const userId = ref("");
const showtimeId = ref("");
const status = ref("");

function handleAdminError(error: unknown) {
  const message = String(error);

  if (message.includes("admin access required")) {
    alert("Admin only");
    router.push("/");
    return true;
  }

  return false;
}

async function loadLogs() {
  try {
    const logRes = await apiGet("/admin/logs");
    logs.value = Array.isArray(logRes) ? logRes : logRes.data || [];
  } catch (error) {
    console.error(error);
    if (handleAdminError(error)) return;
    alert(String(error));
  }
}

async function loadBookings() {
  try {
    loading.value = true;

    const params = new URLSearchParams();

    if (userId.value.trim()) {
      params.append("user_id", userId.value.trim());
    }

    if (showtimeId.value.trim()) {
      params.append("showtime_id", showtimeId.value.trim());
    }

    if (status.value.trim()) {
      params.append("status", status.value.trim());
    }

    const query = params.toString() ? `?${params.toString()}` : "";
    const bookingRes = await apiGet(`/admin/bookings${query}`);

    bookings.value = Array.isArray(bookingRes) ? bookingRes : bookingRes.data || [];
  } catch (error) {
    console.error(error);
    if (handleAdminError(error)) return;
    alert(String(error));
  } finally {
    loading.value = false;
  }
}

async function loadData() {
  await Promise.all([loadBookings(), loadLogs()]);
}

function handleSearch() {
  loadBookings();
}

function clearFilters() {
  userId.value = "";
  showtimeId.value = "";
  status.value = "";
  loadBookings();
}

function formatValue(value: unknown) {
  if (value === null || value === undefined) return "-";
  if (typeof value === "object") return JSON.stringify(value);
  return String(value);
}

onMounted(loadData);
</script>

<template>
  <div class="min-h-screen bg-slate-100 p-6">
    <div class="mx-auto max-w-8xl space-y-6">
      <div
        class="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm md:flex-row md:items-center md:justify-between"
      >
        <div>
          <h1 class="text-3xl font-bold text-slate-900">Admin Dashboard</h1>
          <p class="mt-1 text-sm text-slate-500">
            View all bookings and audit logs
          </p>
        </div>

        <div class="flex gap-3">
          <button
            @click="loadData"
            class="rounded-xl border border-slate-300 bg-white px-4 py-2 font-medium text-slate-700 hover:bg-slate-50"
          >
            Refresh
          </button>
          <button
            @click="router.push('/')"
            class="rounded-xl bg-slate-900 px-4 py-2 font-medium text-white hover:bg-slate-800"
          >
            Back
          </button>
        </div>
      </div>

      <div class="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
        <h2 class="mb-4 text-xl font-semibold text-slate-900">Booking Filters</h2>

        <div class="grid gap-3 md:grid-cols-4">
          <input
            v-model="userId"
            type="text"
            placeholder="User ID"
            class="rounded-xl border border-slate-300 px-4 py-3 outline-none focus:border-slate-500"
          />

          <input
            v-model="showtimeId"
            type="text"
            placeholder="Showtime ID"
            class="rounded-xl border border-slate-300 px-4 py-3 outline-none focus:border-slate-500"
          />

          <select
            v-model="status"
            class="rounded-xl border border-slate-300 px-4 py-3 outline-none focus:border-slate-500"
          >
            <option value="">All Status</option>
            <option value="PENDING">PENDING</option>
            <option value="BOOKED">BOOKED</option>
            <option value="EXPIRED">EXPIRED</option>
          </select>

          <div class="flex gap-2">
            <button
              @click="handleSearch"
              class="flex-1 rounded-xl bg-slate-900 px-4 py-3 font-medium text-white hover:bg-slate-800"
            >
              Search
            </button>
            <button
              @click="clearFilters"
              class="flex-1 rounded-xl border border-slate-300 bg-white px-4 py-3 font-medium text-slate-700 hover:bg-slate-50"
            >
              Clear
            </button>
          </div>
        </div>
      </div>

      <div class="grid gap-6 lg:grid-cols-2">
        <div class="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-slate-900">Bookings</h2>
            <span class="rounded-full bg-slate-100 px-3 py-1 text-sm text-slate-600">
              {{ bookings.length }} items
            </span>
          </div>

          <div v-if="loading" class="py-10 text-center text-slate-500">
            Loading...
          </div>

          <div
            v-else-if="bookings.length === 0"
            class="rounded-xl bg-slate-50 p-6 text-center text-slate-500"
          >
            No bookings found
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full text-sm">
              <thead>
                <tr class="border-b border-slate-200 text-left text-slate-500">
                  <th class="px-3 py-3 font-medium">User</th>
                  <th class="px-3 py-3 font-medium">Seat</th>
                  <th class="px-3 py-3 font-medium">Showtime</th>
                  <th class="px-3 py-3 font-medium">Status</th>
                  <th class="px-3 py-3 font-medium">Created At</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(booking, index) in bookings"
                  :key="booking.id || booking._id || index"
                  class="border-b border-slate-100 last:border-0"
                >
                  <td class="px-3 py-3 text-slate-800">
                    {{ formatValue(booking.user_id || booking.userId) }}
                  </td>
                  <td class="px-3 py-3 text-slate-800">
                    {{ formatValue(booking.seat_number || booking.seatNumber) }}
                  </td>
                  <td class="px-3 py-3 text-slate-800">
                    {{ formatValue(booking.showtime_id || booking.showtimeId) }}
                  </td>
                  <td class="px-3 py-3">
                    <span
                      class="rounded-full px-3 py-1 text-xs font-semibold"
                      :class="{
                        'bg-amber-100 text-amber-700':
                          (booking.status || '').toUpperCase() === 'PENDING',
                        'bg-emerald-100 text-emerald-700':
                          (booking.status || '').toUpperCase() === 'BOOKED',
                        'bg-red-100 text-red-700':
                          (booking.status || '').toUpperCase() === 'EXPIRED',
                        'bg-slate-100 text-slate-700':
                          !['PENDING', 'BOOKED', 'EXPIRED'].includes(
                            (booking.status || '').toUpperCase()
                          ),
                      }"
                    >
                      {{ formatValue(booking.status) }}
                    </span>
                  </td>
                  <td class="px-3 py-3 text-slate-800">
                    {{ formatValue(booking.created_at || booking.createdAt) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-slate-900">Audit Logs</h2>
            <span class="rounded-full bg-slate-100 px-3 py-1 text-sm text-slate-600">
              {{ logs.length }} items
            </span>
          </div>

          <div
            v-if="logs.length === 0"
            class="rounded-xl bg-slate-50 p-6 text-center text-slate-500"
          >
            No logs found
          </div>

          <div v-else class="max-h-[600px] space-y-3 overflow-auto pr-1">
            <div
              v-for="(log, index) in logs"
              :key="log.id || log._id || index"
              class="rounded-xl border border-slate-200 bg-slate-50 p-4"
            >
              <div class="flex flex-wrap items-center gap-2">
                <span class="rounded-full bg-slate-900 px-2.5 py-1 text-xs font-semibold text-white">
                  {{ formatValue(log.event) }}
                </span>
                <span class="text-xs text-slate-500">
                  {{ formatValue(log.created_at) }}
                </span>
              </div>

              <div class="mt-3 grid gap-2 text-sm text-slate-700">
                <p>
                  <span class="font-medium text-slate-900">User:</span>
                  {{ formatValue(log.user_id || log.userId) }}
                </p>
                <p>
                  <span class="font-medium text-slate-900">Seat:</span>
                  {{ formatValue(log.seat_number || log.seatNumber) }}
                </p>
                <p>
                  <span class="font-medium text-slate-900">Showtime:</span>
                  {{ formatValue(log.showtime_id || log.showtimeId) }}
                </p>
                <p>
                  <span class="font-medium text-slate-900">Message:</span>
                  {{ formatValue(log.message || log.description || log.detail) }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>