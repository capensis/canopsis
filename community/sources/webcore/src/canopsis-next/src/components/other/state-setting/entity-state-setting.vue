<template>
  <v-layout column>
    <h5 class="subheading font-weight-bold">
      {{ $t('stateSetting.title') }}
    </h5>
    <v-text-field
      :value="stateSetting.title"
      :loading="stateSettingPending"
      disabled
    />
  </v-layout>
</template>

<script>
import { debounce } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('entity');

export default {
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
    this.debouncedCheckStateSetting = debounce(this.checkStateSetting, 500);

    this.checkStateSetting(this.form);
  },
  methods: {
    ...mapActions({
      checkEntityStateSetting: 'checkStateSetting',
    }),

    async checkStateSetting(form) {
      try {
        this.stateSettingPending = true;
        this.stateSetting = await this.checkEntityStateSetting({ data: this.preparer(form) });
      } catch (err) {
        console.error(err);
      } finally {
        this.stateSettingPending = false;
      }
    },
  },
};
</script>
