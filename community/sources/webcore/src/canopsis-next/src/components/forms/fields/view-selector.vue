<template>
  <v-layout align-center>
    <v-btn
      class="mr-2"
      color="secondary"
      small
      @click="showViewSelectModal"
    >
      {{ $t('user.selectDefaultView') }}
    </v-btn>
    <div>{{ defaultViewTitle }}</div>
    <v-btn
      v-if="value"
      icon
      @click="clearDefaultView"
    >
      <v-icon color="error">
        clear
      </v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';

export default {
  mixins: [entitiesViewGroupMixin],
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
      if (!this.value) {
        return null;
      }

      const userDefaultView = this.getViewById(this.value);

      return userDefaultView ? userDefaultView.title : null;
    },
  },
  methods: {
    showViewSelectModal() {
      this.$modals.show({
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
