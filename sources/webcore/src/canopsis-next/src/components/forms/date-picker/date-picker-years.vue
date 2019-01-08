<script>
import VDatePickerYears from 'vuetify/es5/components/VDatePicker/VDatePickerYears';

export default {
  extends: VDatePickerYears,
  mounted() {
    this.scrollToActive();
  },
  methods: {
    scrollToActive() {
      const activeItem = this.$el.getElementsByClassName('active')[0];

      if (activeItem) {
        this.$el.scrollTop = (activeItem.offsetTop - (this.$el.offsetHeight / 2)) + (activeItem.offsetHeight / 2);
      } else {
        this.$el.scrollTop = (this.$el.scrollHeight / 2) - (this.$el.offsetHeight / 2);
      }
    },
  },
};
</script>
