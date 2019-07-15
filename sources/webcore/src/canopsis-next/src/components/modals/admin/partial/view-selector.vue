<template lang="pug">
  v-layout(align-center)
    v-btn(
    data-test="selectDefaultViewButton",
    color="secondary",
    small,
    @click="showViewSelectModal"
    ) {{ $t('user.selectDefaultView') }}
    div {{ defaultViewTitle }}
    v-btn(v-if="value", icon, @click="clearDefaultView")
      v-icon(color="error") clear
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import entitiesViewMixin from '@/mixins/entities/view';

export default {
  mixins: [modalMixin, entitiesViewMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: null,
    },
  },
  computed: {
    defaultViewTitle() {
      const userDefaultView = this.getViewById(this.value);

      return userDefaultView ? userDefaultView.title : null;
    },
  },
  methods: {
    showViewSelectModal() {
      this.showModal({
        name: MODALS.selectView,
        config: {
          action: viewId => this.$emit('input', viewId),
        },
      });
    },

    clearDefaultView() {
      this.$emit('input', '');
    },
  },
};
</script>
