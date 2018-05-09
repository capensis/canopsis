<template lang="pug">
  v-navigation-drawer(:value="isPanelOpen", disable-resize-watcher, stateless, touchless, right, app)
    v-toolbar(color="blue darken-4")
      v-list
        v-list-tile
          v-list-tile-title(class="title white--text text-xs-center") Alarms list settings
      v-icon(@click.stop="closePanel", color="white") close
    v-divider
    v-list(expand, class="pt-0")
      v-list-group
        v-list-tile(slot="activator", active-class="activeHeader") Title
        v-container
          v-text-field(placeholder="Widget title")
      v-divider
      v-list-group
        v-list-tile(slot="activator") Default Sort column
        v-container
          v-text-field(placeholder="Column name")
          v-select(:items="sortChoices", value="ASC")
      v-divider
      v-list-group
        v-list-tile(slot="activator") Column names
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
                v-text-field(placeholder="Label")
              v-flex(xs11)
                v-text-field(placeholder="Value")
          v-divider
      v-divider
      v-list-group
        v-list-tile(slot="activator") Periodic refresh
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
        v-list-tile(slot="activator") Default number of elements/page
        v-container
          v-text-field(
            placeholder="Elements per page",
            type="number"
          )
      v-divider
      v-list-group
        v-list-tile(slot="activator") Filter on Open/Resolve
        v-container
          v-layout
            v-checkbox(
              label="Open",
              v-model="openCheckbox",
              hide-details
            )
            v-checkbox(
              label="Resolve",
              v-model="resolveCheckbox",
              hide-details
            )
      v-divider
      v-list-group
        v-list-tile(slot="activator") Filters
        v-container
          v-select(label="Select a filter")
      v-divider
      v-list-group(disabled)
        v-list-tile(slot="activator") Info popup
      v-divider
    v-btn(
      color="green darken-4 white--text",
      depressed,
      fixed,
      right
    ) Save
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
</style>
