<template lang="pug">
  v-form.pa-3.pbehavior-form(@submit.prevent="$emit('submit', form)")
    pbehavior-form(v-model="form")
    v-layout(row, justify-end)
      v-btn.mr-0.mb-0(
        depressed,
        flat,
        @click="$emit('close')"
      ) {{ $t('common.cancel') }}
      v-btn.mr-0.mb-0.primary.white--text(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { pbehaviorToComments, pbehaviorToExdates, pbehaviorToForm } from '@/helpers/forms/pbehavior';

import PbehaviorForm from '@/components/other/pbehavior/calendar/partials/pbehavior-form.vue';

export default {
  components: { PbehaviorForm },
  inject: ['$validator'],
  props: {
    calendarEvent: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      form: {
        general: pbehaviorToForm(this.calendarEvent),
        exdate: pbehaviorToExdates(this.calendarEvent),
        comments: pbehaviorToComments(this.calendarEvent),
      },
    };
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-form {
    overflow: auto;
    width: 500px;
    max-height: 600px;
  }
</style>
