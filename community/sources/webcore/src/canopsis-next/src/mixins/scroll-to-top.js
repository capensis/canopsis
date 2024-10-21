import { debounce } from 'lodash';

export const scrollToTopMixin = {
  data() {
    return {
      pageScrolled: false,
    };
  },
  created() {
    this.debouncedCheckScrollPosition = debounce(this.checkScrollPosition, 100);
  },
  mounted() {
    document.addEventListener('scroll', this.debouncedCheckScrollPosition);
  },
  beforeDestroy() {
    document.removeEventListener('scroll', this.debouncedCheckScrollPosition);
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
