import { MODALS } from '@/constants';

import layoutNavigationEditingModeMixin from './editing-mode';

export default {
  mixins: [layoutNavigationEditingModeMixin],
  props: {
    group: {
      type: Object,
      required: true,
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
