<script setup lang="ts">
defineProps<{
  seats: Array<{
    seat_number: string;
    status: string;
  }>;
}>();

const emit = defineEmits<{
  select: [seatNumber: string];
}>();

function seatStyle(status: string) {
  if (status === "BOOKED") return { background: "#ef4444", color: "#fff" };
  if (status === "LOCKED") return { background: "#f59e0b", color: "#fff" };
  return { background: "#22c55e", color: "#fff" };
}
</script>

<template>
  <div
    style="display:grid;grid-template-columns:repeat(5,1fr);gap:12px;max-width:420px"
  >
    <button
      v-for="seat in seats"
      :key="seat.seat_number"
      :style="{
        ...seatStyle(seat.status),
        padding:'16px',
        border:'none',
        borderRadius:'10px',
        cursor: seat.status === 'BOOKED' ? 'not-allowed' : 'pointer',
        opacity: seat.status === 'BOOKED' ? 0.7 : 1,
      }"
      :disabled="seat.status === 'BOOKED'"
      @click="emit('select', seat.seat_number)"
    >
      {{ seat.seat_number }}
      <br />
      {{ seat.status }}
    </button>
  </div>
</template>