<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $t('common.description') }}</v-tab>
    <v-tab-item>
      <v-layout
        class="pa-3 secondary lighten-2"
        column
      >
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <pre>{{ pbehaviorException.description }}</pre>
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
    <v-tab>{{ $t('common.periods') }}</v-tab>
    <v-tab-item>
      <v-layout
        class="pa-3 secondary lighten-2"
        column
      >
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <pbehavior-exceptions-field
                :exdates="preparedPbehaviorExceptionExdates"
                with-exdate-type
                disabled
              />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { computed } from 'vue';

import { pbehaviorExceptionToForm } from '@/helpers/entities/pbehavior/exception/form';

import PbehaviorExceptionsField from '../fields/pbehavior-exceptions-field.vue';

export default {
  components: { PbehaviorExceptionsField },
  props: {
    pbehaviorException: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const preparedPbehaviorExceptionExdates = computed(() => {
      const { exdates } = pbehaviorExceptionToForm(props.pbehaviorException);

      return exdates;
    });

    return {
      preparedPbehaviorExceptionExdates,
    };
  },
};
</script>
