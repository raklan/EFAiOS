class AudioManager {
  constructor() {
    this.sounds = {};
    this.ambientTimer = null;
    this.muted = false;
  }

  load(name, src, volume = 1, loop = false) {
    const audio = new Audio(src);
    audio.volume = volume;
    audio.loop = loop;
    this.sounds[name] = audio;
  }

  play(name) {    
    if(this.muted){
        console.warn("AudioManager.play() was called, but instance is muted. Not playing sound:", name);
        return;
    }

    const sound = this.sounds[name];
    if (!sound){
        console.error('could not find sound:', name);
        return;
    } 
    sound.play();
  }

  stop(name) {
    const sound = this.sounds[name];
    if (!sound) return;
    sound.pause();
    sound.currentTime = 0;
  }

  mute(){
    this.muted = true;
  }

  unmute(){
    this.muted = false;
  }
}