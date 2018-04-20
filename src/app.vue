<template lang="pug">
  v-app#app
    top-bar
    side-bar(
      :isOpen='isSideBarOpen'
    )
    v-content
      nav#nav
        router-link(to="/") Home
        router-link(to="/about") About
      router-view
</template>

<script>
import TopBar from '@/components/layout/top-bar.vue';
import SideBar from '@/components/layout/side-bar.vue';
import { mapState } from 'vuex';

export default {
  name: 'App',
  components: {
    TopBar,
    SideBar,
  },
  data() {
    const locales = [
      { key: 'fr', text: 'Français' },
      { key: 'en', text: 'English' },
    ];

    return {
      currentLocaleIndex: locales.find(({ key }) => key === this.$i18n.locale),
    };
  },
  computed: {
    ...mapState({
      isSideBarOpen: state => state.app.isSideBarOpen,
    }),
  },
};
</script>

