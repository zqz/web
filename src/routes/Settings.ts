const SETTING_THEME = 'theme';

if (!getTheme()) {
  setTheme('light');
}

function setTheme(theme: string) {
  return;
  // console.log('setting theme', theme);
  // localStorage.setItem(SETTING_THEME, theme);
}

function getTheme() {
  return 'light';
  
  // let x = localStorage.getItem(SETTING_THEME);
  // console.log('getting theme', x);
  // return x;
}

const Settings = {
  setTheme,
  getTheme
};

export default Settings;
