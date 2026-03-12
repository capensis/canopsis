<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $t('jobs.tabs.data') }}</v-tab>
    <v-tabs-items mandatory>
      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12 md8 offset-md2>
            <v-card>
              <v-card-text>
                <v-layout class="gap-3" column>
                  <span v-if="item.fail_reason">
                    <strong>{{ $t('jobs.data.webhookFailedPrefix') }} : </strong>
                    {{ item.fail_reason }}
                  </span>
                  <jobs-http-data-collapse-panel
                    v-if="item.raw_request"
                    :title-prefix="$t('jobs.data.request')"
                    :raw="item.raw_request"
                  />
                  <jobs-http-data-collapse-panel
                    v-if="item.raw_response"
                    :title-prefix="$t('jobs.data.response')"
                    :raw="item.raw_response"
                  />
                </v-layout>
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import JobsHttpDataCollapsePanel from './jobs-http-data-collapse-panel.vue';

export default {
  components: { JobsHttpDataCollapsePanel },
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
};
</script>
