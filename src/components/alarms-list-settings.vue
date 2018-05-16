<template lang="pug">
  v-navigation-drawer(:value="isPanelOpen", disable-resize-watcher, stateless, touchless, right, app)
    v-toolbar(color="blue darken-4")
      v-list
        v-list-tile
          v-list-tile-title(class="white--text text-xs-center") {{$t('alarm_list_settings.alarm_list_settings')}}
      v-icon(@click.stop="closePanel", color="white" class="closeIcon") close
    v-divider
    v-list(expand, class="pt-0")
      v-list-group
        v-list-tile(slot="activator", active-class="activeHeader") {{$t('common.title')}}
        v-container
          v-text-field(:placeholder="$t('alarm_list_settings.widget_title')")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('alarm_list_settings.default_sort_column')}}
        v-container
          v-text-field(:placeholder="$t('alarm_list_settings.column_name')")
          v-select(:items="sortChoices", value="ASC")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('alarm_list_settings.column_names')}}
        v-container
          v-card
            v-layout(justify-space-between class="pt-2")
              v-flex(xs3)
                v-layout(justify-space-between class="text-xs-center pl-2")
                  v-flex(xs1)
                    v-icon arrow_upward
                  v-flex(xs5)
                    v-icon arrow_downward
              v-flex(xs3 class="d-flex")
                div(class="text-xs-right pr-2")
                  v-icon(color="red") close
            v-layout(justify-center wrap)
              v-flex(xs11)
                v-text-field(:placeholder="$t('common.label')")
              v-flex(xs11)
                v-text-field(:placeholder="$t('common.value')")
          v-divider
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('alarm_list_settings.periodic_refresh')}}
        v-container
          v-layout
            v-flex
              v-switch(
                v-model="isPeriodicRefreshEnable",
                color="green darken-3",
                hide-details,
              )
            v-flex
              v-text-field(
                class='pt-0',
                hide-details,
                type="number",
                :disabled="!isPeriodicRefreshEnable")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('alarm_list_settings.default_number_of_elements_per_page')}}
        v-container
          v-text-field(
            :placeholder="$t('alarm_list_settings.elements_per_page')",
            type="number"
          )
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('alarm_list_settings.filter_on_open_resolved')}}
        v-container
          v-layout
            v-checkbox(
              :label="$t('alarm_list_settings.open')",
              v-model="openCheckbox",
              hide-details
            )
            v-checkbox(
              :label="$t('alarm_list_settings.resolved')",
              v-model="resolveCheckbox",
              hide-details
            )
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{$t('alarm_list_settings.filters')}}
        v-container
          v-select(:label="$t('alarm_list_settings.select_a_filter')")
      v-divider
      v-list-group(disabled)
        v-list-tile(slot="activator") {{$t('alarm_list_settings.info_popup')}}
      v-divider
    v-btn(
      color="green darken-4 white--text",
      depressed,
      fixed,
      right
    ) {{$t('common.save')}}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('AlarmsListSettings');

export default {
  name: 'AlarmsListSettings',
  data() {
    return {
      sortChoices: ['ASC', 'DESC'],
      isPeriodicRefreshEnable: false,
      openCheckbox: true,
      resolveCheckbox: false,
    };
  },
  computed: {
    ...mapGetters(['isPanelOpen']),
  },
  methods: {
    ...mapActions(['closePanel']),
  },
};
</script>

<style scoped>
  .activeHeader {
    background-color: blue;
  }

  .closeIcon:hover {
    cursor: pointer;
  }
</style>
