<template>
  <div class="d-flex flex-column gap-2 mt-3">
    <div class="d-flex align-center">
      <span class="text-body-2 font-weight-bold">
        {{ t('settings.associateTicket.ticketFormats') }}
      </span>
      <c-help-icon
        :text="t('settings.associateTicket.formatsHelp')"
        icon-class="ml-2"
        max-width="400"
        top
      />
    </div>
    <span>{{ t('settings.associateTicket.description') }}</span>

    <ticket-formats-form v-model="formats" />
  </div>
</template>

<script>
import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import TicketFormatsForm from '@/components/forms/ticket-formats-form.vue';

export default {
  components: { TicketFormatsForm },

  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const formats = form.model('formats');

    return {
      t,
      formats,
    };
  },
};
</script>
