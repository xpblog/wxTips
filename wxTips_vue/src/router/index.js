import { createRouter, createWebHistory } from 'vue-router'
import Qr from '../views/Qr.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Qr',
      component: Qr
    }
  ]
})

export default router
