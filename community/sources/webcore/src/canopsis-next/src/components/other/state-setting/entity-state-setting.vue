<template>
  <v-layout column>
    <h5 class="subheading font-weight-bold">
      {{ $t('stateSetting.title') }}
    </h5>
    <v-text-field
      :value="stateSetting?.title"
      :loading="stateSettingPending"
      disabled
    />
  </v-layout>
</template>

<script>
import { debounce } from 'lodash';

import { checkStateSettingMixin } from '@/mixins/entities/check-state-setting';

export default {
  mixins: [checkStateSettingMixin],
  props: {
    form: {
      type: Object,
      required: true,
    },
    preparer: {
      type: Function,
      default: () => d => d,
    },
  },
  data() {
    return {
      stateSetting: {},
      stateSettingPending: false,
    };
  },
  watch: {
    form: {
      deep: true,
      handler(form) {
        this.debouncedCheckStateSetting(form);
      },
    },
  },
  created() {
    this.debouncedCheckStateSetting = debounce(this.checkStateSettingByForm, 500);

    this.checkStateSettingByForm(this.form);
  },
  methods: {
    checkStateSettingByForm(form) {
      return this.checkStateSetting(this.preparer(form));
    },
  },
};
</script>
