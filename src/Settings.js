const SETTING_THEME = 'theme';

if (!getTheme()) {
  setTheme('light');
}

function setTheme(theme) {
  console.log('setting theme', theme);
  localStorage.setItem(SETTING_THEME, theme);
}

function getTheme() {
  let x = localStorage.getItem(SETTING_THEME);
  console.log('getting theme', x);
  return x;
}

const Settings = {
  setTheme,
  getTheme
};

export default Settings;
