import { useI18n } from '@/hooks/i18n'

export const useHealthcheckNodes = () => {
  const { t, te } = useI18n()

  const getNodeName = (nodeName) => {
    const nodeLabelKey = `healthcheck.nodes.${nodeName}.name`

    return te(nodeLabelKey) ? t(nodeLabelKey) : nodeName
  }

  const getTooltipText = (data = {}) => {
    const statusKeys = []
    const itemWithDefaultFlags = {
      is_running: true,
      is_queue_overflown: false,
      is_too_few_instances: false,
      is_diff_instances_config: false,
      ...data,
    }

    if (!itemWithDefaultFlags.is_running) {
      statusKeys.push('healthcheck.notRunning')
    }

    if (itemWithDefaultFlags.is_queue_overflown) {
      statusKeys.push('healthcheck.queueOverflow')
    }

    if (itemWithDefaultFlags.is_too_few_instances) {
      statusKeys.push('healthcheck.lackOfInstances')
    }

    if (itemWithDefaultFlags.is_diff_instances_config) {
      statusKeys.push('healthcheck.diffInstancesConfig')
    }

    if (itemWithDefaultFlags.is_unknown) {
      statusKeys.push('healthcheck.unknown')
    }

    return statusKeys
      .map(stateKey => t(stateKey, { name: getNodeName(data.name) }))
      .join('\n')
  }

  return {
    getNodeName,
    getTooltipText,
  }
}
