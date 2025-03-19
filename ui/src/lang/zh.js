export default {
  common: {
    text: {
      login: '登 录',
      loging: '登 录 中...',
      remember: '记住密码',
    }
  },
  navbar: {
    text: {
      vaults: '保险库',
      setting: '系统设置',
      logout: '退出登录'
    }
  },
  vaults: {
    text: {
      hosts: '主机',
      port_forwarding: '端口转发',
      history: '历史记录',
    }
  },
  hosts: {
    text: {
      new_host: '新建主机',
      hosts_title: '主机列表',
      add_host_title: '新建主机',
      edit_host_title: '主机详情',
      host_general: '通用设置',
      host_label: '标签',
      host_address: '地址',
      host_ip: 'IP 或主机名',
      login_by_password: '凭据（通过密码）',
      username: '用户名',
      password: '密码',
      connect: '连接',
      remove: '删除'
    }
  },

  forwarding: {
    text: {
      new_host: '新建端口转发',
      hosts_title: '端口转发列表',
      add_host_title: '新建本地端口转发',
      edit_host_title: '本地端口转发详情',
      host_general: '通用设置',
      host_label: '标签',
      host_address_local: '本地地址',
      local_address: '绑定地址',
      local_port: '绑定端口',
      intermediate: '跳板服务器',
      host_address_remote: '远端地址',
      remote_address: '远端地址',
      remote_port: '远端端口',

      connect: '连接',
      remove: '删除',


      select_host: '选择主机',
    }
  },
  setting: {
    text: {
      setting: '设置',
      account: '设置账号',
      passwordV: '修改密码',
      theme: '主题',
      email: {
        title: '邮箱',
        placeholder: '请输入您的邮箱',
        submit_button: '保存',

        save_success: '邮箱修改成功',
        save_error: '修改邮箱失败'
      },
      password: {
        old_password: '当前密码',
        old_placeholder: '请输入当前密码',
        new_password: '新密码',
        new_placeholder: '请输入新密码',
        confirm_password: '密码确认',
        confirm_placeholder: '请确认您的密码',
        submit_button: '保存',

        save_success: '密码修改成功',
        save_error: '当前密码输入错误'
      }
    }
  }
}