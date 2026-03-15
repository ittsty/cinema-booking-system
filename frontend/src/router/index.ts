import { createRouter, createWebHistory } from "vue-router";
import { onAuthStateChanged } from "firebase/auth";
import { auth } from "../firebase/config";
import { apiGet } from "../api/client";

import LoginView from "../views/LoginView.vue";
import HomeView from "../views/HomeView.vue";
import SeatMapView from "../views/SeatMapView.vue";
import AdminView from "../views/AdminView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/login", component: LoginView },
    { path: "/", component: HomeView, meta: { requiresAuth: true } },
    { path: "/showtimes/:id/seats", component: SeatMapView, meta: { requiresAuth: true } },
    { path: "/admin", component: AdminView, meta: { requiresAuth: true, requiresAdmin: true } },
  ],
});

function waitForAuth() {
  return new Promise<void>((resolve) => {
    const unsub = onAuthStateChanged(auth, () => {
      unsub();
      resolve();
    });
  });
}

router.beforeEach(async (to) => {
  await waitForAuth();

  if (to.meta.requiresAuth && !auth.currentUser) {
    return "/login";
  }
  if (to.path === "/login" && auth.currentUser) {
    return "/";
  }
  if (to.meta.requiresAdmin) {
    try {
      const me = await apiGet("/me");

      if (me.role !== "ADMIN") {
        alert("Admin only");
        return "/";
      }
    } catch (error) {
      console.error(error);
      return "/";
    }
  }
});

export default router;