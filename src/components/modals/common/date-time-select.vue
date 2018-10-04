<template lang="pug">
  v-card
    v-card-title.green.darken-4.white--text
      v-layout(justify-space-between, align-center)
        h2 {{ $t('modals.dateTimeSelect.title') }}
        v-btn(@click='hideModal', icon, small)
          v-icon.white--text close
    v-card-text
      date-time-picker(v-model="date")
      v-btn(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import DateTimePicker from '@/components/forms/date-time-picker.vue';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

export default {
  name: MODALS.dateSelect,
  components: {
    DateTimePicker,
  },
  mixins: [modalInnerMixin],
  data() {
    const date = this.modal.config.value || new Date();
    return {
      date,
    };
  },
  methods: {
    submit() {
      if (this.config.action) {
        this.config.action(this.date);
      }
      this.hideModal();
    },
  },
};
</script>

