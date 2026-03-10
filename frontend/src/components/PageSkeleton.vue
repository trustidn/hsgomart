<template>
  <div class="p-4 lg:p-6 space-y-6">
    <!-- Title bar -->
    <div class="flex items-center justify-between gap-4">
      <div class="skeleton h-6 w-48 rounded-lg" />
      <div class="skeleton h-9 w-24 rounded-lg" />
    </div>

    <!-- Cards row (Dashboard-style) -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div v-for="i in 4" :key="i" class="rounded-xl border border-gray-200 dark:border-gray-800 overflow-hidden bg-white dark:bg-gray-900">
        <div class="p-4 space-y-3">
          <div class="skeleton h-4 w-20 rounded" />
          <div class="skeleton h-8 w-16 rounded" />
          <div class="skeleton h-3 w-full rounded" />
        </div>
      </div>
    </div>

    <!-- Content block (table/list style) -->
    <div class="rounded-xl border border-gray-200 dark:border-gray-800 overflow-hidden bg-white dark:bg-gray-900">
      <div class="p-4 border-b border-gray-100 dark:border-gray-800 flex flex-wrap gap-2">
        <div class="skeleton h-10 w-full max-w-[200px] rounded-lg" />
        <div class="skeleton h-10 w-28 rounded-lg" />
      </div>
      <div class="divide-y divide-gray-100 dark:divide-gray-800">
        <div v-for="i in 6" :key="i" class="px-4 py-4 flex items-center gap-4">
          <div class="skeleton h-10 w-10 rounded-lg shrink-0" />
          <div class="flex-1 min-w-0 space-y-2">
            <div class="skeleton h-4 w-48 rounded" />
            <div class="skeleton h-3 w-32 rounded" />
          </div>
          <div class="skeleton h-8 w-16 rounded shrink-0" />
          <div class="skeleton h-8 w-20 rounded shrink-0" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.skeleton {
  background-color: #e5e7eb;
  background-image: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.15) 50%,
    transparent 100%
  );
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
}

.dark .skeleton {
  background-color: #374151;
  background-image: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.06) 50%,
    transparent 100%
  );
}

@keyframes skeleton-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
