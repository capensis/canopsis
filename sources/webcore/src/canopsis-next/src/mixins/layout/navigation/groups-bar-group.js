import { MODALS } from '@/constants';

export default {
  props: {
    group: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    showEditGroupModal() {
      this.$modals.show({
        name: MODALS.createGroup,
        config: {
          group: this.group,
        },
      });
    },
  },
};
