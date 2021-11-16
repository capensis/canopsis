export const scrollToTopMixin = {
  data() {
    return {
      pageScrolled: false,
    };
  },

  mounted() {
    document.addEventListener('scroll', this.checkScrollPosition);
  },

  beforeDestroy() {
    document.removeEventListener('scroll', this.checkScrollPosition);
  },
  methods: {
    checkScrollPosition() {
      this.pageScrolled = window.scrollY > 0;
    },

    scrollToTop() {
      window.scrollTo({
        top: 0,
        behavior: 'smooth',
      });
    },
  },
};
