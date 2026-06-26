import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { API_ROUTES } from '@/config'
import { POPUP_TYPES } from '@/constants'

import request from '@/services/request'

import { durationToSeconds } from '@/helpers/date/duration'
import { getLocalTimezone } from '@/helpers/date/date'

export const useAppInfoStore = defineStore('appInfo', () => {
  const appInfo = ref({})
  const pending = ref(false)

  const version = computed(() => appInfo.value.version)
  const logo = computed(() => appInfo.value.logo)
  const appTitle = computed(() => appInfo.value.app_title)
  const maintenance = computed(() => appInfo.value.maintenance)
  const defaultColorTheme = computed(() => appInfo.value.default_color_theme)
  const popupTimeout = computed(() => appInfo.value.popup_timeout || {})
  const maxMatchedItems = computed(() => appInfo.value.max_matched_items)
  const checkCountRequestTimeout = computed(() => appInfo.value.check_count_request_timeout)
  const allowChangeSeverityToInfo = computed(() => appInfo.value.allow_change_severity_to_info)
  const showHeaderOnKioskMode = computed(() => appInfo.value.show_header_on_kiosk_mode)
  const requiredInstructionApprove = computed(() => appInfo.value.required_instruction_approve)
  const footer = computed(() => appInfo.value.footer)
  const edition = computed(() => appInfo.value.edition)
  const description = computed(() => appInfo.value.login_page_description)
  const language = computed(() => appInfo.value.language)
  const timezone = computed(() => appInfo.value.timezone)
  const fileUploadMaxSize = computed(() => appInfo.value.file_upload_max_size ?? 0)
  const fileImportMaxSize = computed(() => appInfo.value.file_import_max_size ?? 0)
  const remediationJobConfigTypes = computed(() => appInfo.value.remediation?.job_config_types ?? [])
  const casConfig = computed(() => appInfo.value.login?.casconfig)
  const samlConfig = computed(() => appInfo.value.login?.saml2config)
  const oauthConfig = computed(() => appInfo.value.login?.oauth2config)

  const isBasicAuthEnabled = computed(() => !!appInfo.value.login?.basic?.enable)
  const isLDAPAuthEnabled = computed(() => !!appInfo.value.login?.ldapconfig?.enable)
  const isCASAuthEnabled = computed(() => !!appInfo.value.login?.casconfig?.enable)
  const isSAMLAuthEnabled = computed(() => !!appInfo.value.login?.saml2config?.enable)
  const isOauthAuthEnabled = computed(() => !!appInfo.value.login?.oauth2config?.enable)

  const eventsCountTriggerDefaultThreshold = computed(() => (
    appInfo.value.events_count_trigger_default_threshold
  ))

  const disabledTransitions = computed(() => appInfo.value.disabled_transitions)
  const autoSuggestPbehaviorName = computed(() => appInfo.value.auto_suggest_pbehavior_name)
  const userTimezones = computed(() => appInfo.value.user_timezones ?? [])

  const shownUserTimezone = computed(() => (
    userTimezones.value.length > 1
    || (userTimezones.value.length && !userTimezones.value.includes(getLocalTimezone()))
  ))

  const serialName = computed(() => appInfo.value.serial_name)
  const versionUpdated = computed(() => appInfo.value.version_updated)
  const versionDescription = computed(() => appInfo.value.version_description ?? '')
  const notificationDisplayCount = computed(() => appInfo.value.notification_display_count ?? 0)

  const fetchAppInfo = async ({ setGlobalLocale, setDefaultCloseTime } = {}) => {
    try {
      pending.value = true

      const nextAppInfo = await request.get(API_ROUTES.infos.app)

      appInfo.value = nextAppInfo
      pending.value = false

      if (nextAppInfo.language) {
        await setGlobalLocale?.(nextAppInfo.language)
      }

      if (nextAppInfo.popup_timeout) {
        setPopupTimeouts({
          popupTimeout: nextAppInfo.popup_timeout,
          setDefaultCloseTime,
        })
      }
    } catch (err) {
      pending.value = false

      throw err
    }
  }

  const updateUserInterface = ({ data } = {}) => (
    request.post(API_ROUTES.infos.userInterface, data)
  )

  const updateMaintenanceMode = ({ data } = {}) => (
    request.put(API_ROUTES.maintenance, data)
  )

  const setPopupTimeouts = ({ popupTimeout: nextPopupTimeout = {}, setDefaultCloseTime } = {}) => {
    const timeInfo = durationToSeconds(nextPopupTimeout.info) * 1000
    const timeError = durationToSeconds(nextPopupTimeout.error) * 1000

    setDefaultCloseTime?.({ type: POPUP_TYPES.success, time: timeInfo })
    setDefaultCloseTime?.({ type: POPUP_TYPES.info, time: timeInfo })
    setDefaultCloseTime?.({ type: POPUP_TYPES.warning, time: timeError })
    setDefaultCloseTime?.({ type: POPUP_TYPES.error, time: timeError })
  }

  return {
    appInfo,
    pending,

    version,
    logo,
    appTitle,
    maintenance,
    defaultColorTheme,
    popupTimeout,
    maxMatchedItems,
    checkCountRequestTimeout,
    allowChangeSeverityToInfo,
    showHeaderOnKioskMode,
    requiredInstructionApprove,
    footer,
    edition,
    description,
    language,
    timezone,
    fileUploadMaxSize,
    fileImportMaxSize,
    remediationJobConfigTypes,
    casConfig,
    samlConfig,
    oauthConfig,
    isBasicAuthEnabled,
    isLDAPAuthEnabled,
    isCASAuthEnabled,
    isSAMLAuthEnabled,
    isOauthAuthEnabled,
    eventsCountTriggerDefaultThreshold,
    disabledTransitions,
    autoSuggestPbehaviorName,
    userTimezones,
    shownUserTimezone,
    serialName,
    versionUpdated,
    versionDescription,
    notificationDisplayCount,

    fetchAppInfo,
    updateUserInterface,
    updateMaintenanceMode,
    setPopupTimeouts,
  }
})
