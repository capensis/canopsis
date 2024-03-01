<template>
  <v-layout column>
    <field-title v-field="form.title" />
    <field-periodic-refresh v-field="form.parameters" />
  </v-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';

export default {
  components: {
    FieldTitle,
    FieldPeriodicRefresh,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
