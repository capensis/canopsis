<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $tc('common.information') }}</v-tab>
    <v-tabs-items mandatory>
      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card class="pa-3">
              <v-layout>
                <v-flex
                  offset-xs1
                  xs10
                >
                  <state-setting-information :state-setting="stateSetting" />
                </v-flex>
              </v-layout>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import StateSettingInformation from './state-setting-information.vue';

export default {
  components: { StateSettingInformation },
  props: {
    stateSetting: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
