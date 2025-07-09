<template>
  <v-menu
    :close-on-content-click="false"
    :min-width="400"
    top
    offset-y
  >
    <template #activator="{ on }">
      <v-badge color="blue-grey lighten-1" overlap>
        <v-btn
          color="grey darken-2"
          fab
          dark
          v-on="on"
        >
          <v-icon>$vuetify.icons.manual_instruction</v-icon>
        </v-btn>
        <template #badge>
          <span>9+</span>
        </template>
      </v-badge>
    </template>
    <v-card>
      <v-card-text>
        <v-tabs fixed-tabs>
          <v-tab>Manual</v-tab>
          <v-tab>Manual simplified</v-tab>
          <v-tab-item>
            <v-list>
              <v-list-item>
                <v-list-item-content>
                  <v-list-item-title class="primary--text cursor-pointer mb-3">
                    Instruction name
                  </v-list-item-title>
                  <v-list-item-subtitle>
                    <v-layout class="gap-2" column>
                      <span>Alarm display name</span>
                      <span><strong>Jobs: </strong><v-chip color="warning" small>3 in progress</v-chip></span>
                      <span><strong>Step 1a: </strong>Step name</span>
                      <span class="mt-1 grey--text">
                        <v-icon color="grey" small>timer</v-icon>
                        <span class="ml-2">9 minutes to complete</span>
                      </span>
                    </v-layout>
                  </v-list-item-subtitle>
                </v-list-item-content>
                <v-list-item-action>
                  <c-action-btn
                    type="delete"
                    icon="close"
                    small
                    icon-small
                  />
                </v-list-item-action>
              </v-list-item>
              <v-divider />
            </v-list>
          </v-tab-item>
        </v-tabs>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script>
export default {

};
</script>
