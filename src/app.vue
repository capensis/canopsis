<template lang="pug">
  v-app#app
    #nav
      router-link(to="/") Home
      router-link(to="/about") About
    router-view
</template>
