<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, computed } from "vue";
import { signOut } from "firebase/auth";
import { useRouter, useRoute } from "vue-router";
import { auth } from "../firebase/config";
import { apiGet, apiPost } from "../api/client";

type Seat = {
  seat_number: string;
  status: string;
};

const router = useRouter();
const route = useRoute();
const showtimeId = computed(() => String(route.params.id));
const seats = ref<Seat[]>([]);
const selectedSeat = ref("");
const loading = ref(false);
let socket: WebSocket | null = null;

const selectedSeatStatus = computed(() => {
  return (
    seats.value.find((seat) => seat.seat_number === selectedSeat.value)
      ?.status || ""
  );
});

async function loadSeats() {
  seats.value = await apiGet(`/showtimes/${showtimeId.value}/seats`);
}

function selectSeat(seatNumber: string) {
  selectedSeat.value = seatNumber;
}
async function lockAndCreateBooking() {
  if (!selectedSeat.value) return;

  try {
    loading.value = true;

    await apiPost(
      `/seats/${selectedSeat.value}/lock?showtime_id=${showtimeId.value}`
    );

    await apiPost(`/booking`, {
      seat_number: selectedSeat.value,
      showtime_id: showtimeId.value,
    });

    await loadSeats();
    alert("Seat locked and booking created");
  } catch (error) {
    alert(String(error));
  } finally {
    loading.value = false;
  }
}
// async function lockSeat() {
//   if (!selectedSeat.value) return;
//   try {
//     loading.value = true;
//     await apiPost(
//       `/seats/${selectedSeat.value}/lock?showtime_id=${showtimeId.value}`,
//     );
//     await loadSeats();
//   } catch (error) {
//     alert(String(error));
//   } finally {
//     loading.value = false;
//   }
// }

// async function createBooking() {
//   if (!selectedSeat.value) return;
//   try {
//     loading.value = true;
//     await apiPost(`/booking`, {
//       seat_number: selectedSeat.value,
//       showtime_id: showtimeId.value,
//     });
//     alert("Booking created");
//   } catch (error) {
//     alert(String(error));
//   } finally {
//     loading.value = false;
//   }
// }

async function confirmBooking() {
  if (!selectedSeat.value) return;
  try {
    loading.value = true;
    await apiPost(`/booking/${selectedSeat.value}/confirm`, {
      showtime_id: showtimeId.value,
    });
    await loadSeats();
    alert("Booking confirmed");
  } catch (error) {
    alert(String(error));
  } finally {
    loading.value = false;
  }
}

async function logout() {
  await signOut(auth);
  router.push("/login");
}

function seatClasses(status: string, seatNumber: string) {
  const isSelected = selectedSeat.value === seatNumber;

  if (status === "BOOKED") {
    return [
      "bg-red-500 text-white hover:bg-red-500 cursor-not-allowed",
      isSelected ? "ring-4 ring-red-200" : "",
    ].join(" ");
  }

  if (status === "LOCKED") {
    return [
      "bg-amber-500 text-white hover:bg-amber-500",
      isSelected ? "ring-4 ring-amber-200" : "",
    ].join(" ");
  }

  return [
    "bg-emerald-500 text-white hover:bg-emerald-600",
    isSelected ? "ring-4 ring-emerald-200" : "",
  ].join(" ");
}

function connectWs() {
  socket = new WebSocket("ws://localhost:8080/ws");

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);

      if (
        data.event === "seat_updated" &&
        data.showtime_id === showtimeId.value
      ) {
        const target = seats.value.find(
          (s) => s.seat_number === data.seat_number,
        );
        if (target) {
          target.status = data.status;
        }
      }
    } catch (error) {
      console.error(error);
    }
  };
}

onMounted(async () => {
  await loadSeats();
  connectWs();
});

onBeforeUnmount(() => {
  socket?.close();
});
</script>

<template>
  <div class="min-h-screen bg-slate-100 p-6">
    <div class="mx-auto max-w-7xl space-y-6">
      <div
        class="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm md:flex-row md:items-center md:justify-between"
      >
        <div>
          <h1 class="text-3xl font-bold text-slate-900">Seat Map</h1>
          <p class="mt-1 text-sm text-slate-500">
            Showtime ID: {{ showtimeId }}
          </p>
        </div>

        <div class="flex flex-wrap gap-3">
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

      <div class="grid gap-6 lg:grid-cols-[1.4fr_0.8fr]">
        <div class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
          <div class="mb-5 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-slate-900">
              Choose your seat
            </h2>
            <div class="flex flex-wrap gap-3 text-sm">
              <div class="flex items-center gap-2">
                <span class="h-3 w-3 rounded-full bg-emerald-500"></span>
                <span class="text-slate-600">AVAILABLE</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="h-3 w-3 rounded-full bg-amber-500"></span>
                <span class="text-slate-600">LOCKED</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="h-3 w-3 rounded-full bg-red-500"></span>
                <span class="text-slate-600">BOOKED</span>
              </div>
            </div>
          </div>

          <div
            class="mb-6 rounded-2xl bg-slate-900 py-3 text-center text-sm font-semibold tracking-[0.3em] text-white"
          >
            SCREEN
          </div>

          <div
            class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5"
          >
            <button
              v-for="seat in seats"
              :key="seat.seat_number"
              @click="selectSeat(seat.seat_number)"
              :disabled="seat.status === 'BOOKED'"
              class="rounded-2xl px-4 py-5 text-center font-semibold shadow-sm transition disabled:opacity-80"
              :class="seatClasses(seat.status, seat.seat_number)"
            >
              <div class="text-base">{{ seat.seat_number }}</div>
              <div class="mt-1 text-xs font-medium opacity-90">
                {{ seat.status }}
              </div>
            </button>
          </div>
        </div>

        <div class="space-y-6">
          <div
            class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"
          >
            <h2 class="text-xl font-semibold text-slate-900">Selection</h2>

            <div class="mt-4 space-y-3 text-sm">
              <div class="rounded-xl bg-slate-50 p-4">
                <p class="text-slate-500">Selected Seat</p>
                <p class="mt-1 text-lg font-semibold text-slate-900">
                  {{ selectedSeat || "-" }}
                </p>
              </div>

              <div class="rounded-xl bg-slate-50 p-4">
                <p class="text-slate-500">Current Status</p>
                <p class="mt-1 text-lg font-semibold text-slate-900">
                  {{ selectedSeatStatus || "-" }}
                </p>
              </div>
            </div>

            <div class="mt-5 grid grid-cols-1 gap-3">
              <button
                @click="lockAndCreateBooking"
                :disabled="!selectedSeat || loading"
                class="rounded-xl bg-blue-600 px-4 py-3 font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
              >
                Create Booking
              </button>

              <button
                @click="confirmBooking"
                :disabled="!selectedSeat || loading"
                class="rounded-xl bg-emerald-600 px-4 py-3 font-medium text-white hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
              >
                Confirm Booking
              </button>
            </div>
          </div>

          <div
            class="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"
          >
            <h2 class="text-xl font-semibold text-slate-900">How it works</h2>
            <div class="mt-4 space-y-3 text-sm text-slate-600">
              <p>1. Select a seat from the grid.</p>
              <p>2. Create a booking for the locked seat.</p>
              <p>3. Confirm the booking to complete the purchase flow.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
