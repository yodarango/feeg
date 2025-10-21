// DOM Elements
const bgElement = document.getElementById("background");
const bgVideo = document.getElementById("backgroundVideo");
const bgButton = document.getElementById("bgButton");
const bgModal = document.getElementById("bgModal");
const closeModal = document.getElementById("closeModal");

const soundButton = document.getElementById("soundButton");
const soundModal = document.getElementById("soundModal");
const closeSoundModal = document.getElementById("closeSoundModal");
const globalVolumeSlider = document.getElementById("globalVolumeSlider");
const globalVolumeValue = document.getElementById("globalVolumeValue");
const activeSoundsContainer = document.getElementById("activeSoundsContainer");
const activeSoundsList = document.getElementById("activeSoundsList");
const soundButtonsGrid = document.querySelectorAll(".sound-button-grid");

// State management
const audioPlayers = new Map();
const soundSettings = new Map();
const activeSounds = new Set();
let globalVolume = 0.5;
let expandedSound = null;

// Initialize sound settings structure
function initSoundSettings(soundName) {
  if (!soundSettings.has(soundName)) {
    soundSettings.set(soundName, {
      volume: 0.5,
      speed: 1.0,
      loopTime: 0,
    });
  }
}

// ===== BACKGROUND MODAL =====
bgButton.addEventListener("click", () => {
  bgModal.classList.add("active");
});

closeModal.addEventListener("click", () => {
  bgModal.classList.remove("active");
});

bgModal.addEventListener("click", (e) => {
  if (e.target === bgModal) {
    bgModal.classList.remove("active");
  }
});

// Handle background grid item selection
document.querySelectorAll(".bg-grid-item").forEach((item) => {
  item.addEventListener("click", () => {
    const filename = item.dataset.bgPath;
    const type = item.dataset.bgType;

    // Update active state
    document.querySelectorAll(".bg-grid-item").forEach((i) => {
      i.classList.remove("active");
    });
    item.classList.add("active");

    if (filename) {
      if (type === "video") {
        bgElement.style.backgroundImage = "none";
        bgVideo.src = `/public/bkgs/${filename}`;
        bgVideo.style.display = "block";
      } else {
        bgVideo.style.display = "none";
        bgElement.style.backgroundImage = `url('/public/bkgs/${filename}')`;
      }

      fetch("/api/save-background", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ background: filename }),
      }).catch((error) => console.error("Failed to save background:", error));

      bgModal.classList.remove("active");
    }
  });
});

// Restore previously selected background
fetch("/api/get-background")
  .then((response) => response.json())
  .then((data) => {
    if (data.background) {
      const item = document.querySelector(
        `.bg-grid-item[data-bg-path="${data.background}"]`
      );
      if (item) {
        item.click();
      }
    }
  })
  .catch((error) => console.error("Failed to load background:", error));

// ===== SOUND MODAL =====
soundButton.addEventListener("click", () => {
  soundModal.classList.add("active");
});

closeSoundModal.addEventListener("click", () => {
  soundModal.classList.remove("active");
});

soundModal.addEventListener("click", (e) => {
  if (e.target === soundModal) {
    soundModal.classList.remove("active");
  }
});

// Handle global volume slider
globalVolumeSlider.addEventListener("input", (e) => {
  globalVolume = e.target.value / 100;
  globalVolumeValue.textContent = e.target.value + "%";

  audioPlayers.forEach((audio, soundName) => {
    if (!audio.paused) {
      const settings = soundSettings.get(soundName);
      if (settings) {
        audio.volume = globalVolume * settings.volume;
      }
    }
  });
});

// Function to update the active sounds display
function updateActiveSoundsDisplay() {
  activeSoundsList.innerHTML = "";

  if (activeSounds.size === 0) {
    activeSoundsContainer.style.display = "none";
    return;
  }

  activeSoundsContainer.style.display = "block";

  activeSounds.forEach((soundName) => {
    const settings = soundSettings.get(soundName);
    const soundCard = document.createElement("div");
    soundCard.className = "sound-card";
    soundCard.dataset.soundName = soundName;

    // Compact view (default)
    const compactView = document.createElement("div");
    compactView.className = "sound-card-compact";

    const soundTitle = document.createElement("div");
    soundTitle.className = "sound-card-title";
    soundTitle.textContent = soundName;

    const cardActions = document.createElement("div");
    cardActions.className = "sound-card-actions";

    const editButton = document.createElement("button");
    editButton.className = "sound-card-edit";
    editButton.innerHTML = '<ion-icon name="pencil"></ion-icon>';
    editButton.title = "Edit settings";

    const removeButton = document.createElement("button");
    removeButton.className = "sound-card-remove";
    removeButton.innerHTML = '<ion-icon name="close"></ion-icon>';
    removeButton.title = "Remove sound";

    cardActions.appendChild(editButton);
    cardActions.appendChild(removeButton);

    compactView.appendChild(soundTitle);
    compactView.appendChild(cardActions);
    soundCard.appendChild(compactView);

    // Expanded view (hidden by default)
    const expandedView = document.createElement("div");
    expandedView.className = "sound-card-expanded";
    expandedView.style.display = "none";

    // Volume control
    const volumeSection = document.createElement("div");
    volumeSection.className = "sound-setting";

    const volumeLabel = document.createElement("label");
    volumeLabel.textContent = "Volume";

    const volumeSlider = document.createElement("input");
    volumeSlider.type = "range";
    volumeSlider.className = "volume-slider";
    volumeSlider.min = "0";
    volumeSlider.max = "100";
    volumeSlider.value = settings.volume * 100;

    const volumeValue = document.createElement("span");
    volumeValue.className = "setting-value";
    volumeValue.textContent = Math.round(settings.volume * 100) + "%";

    volumeSection.appendChild(volumeLabel);
    volumeSection.appendChild(volumeSlider);
    volumeSection.appendChild(volumeValue);

    // Speed control
    const speedSection = document.createElement("div");
    speedSection.className = "sound-setting";

    const speedLabel = document.createElement("label");
    speedLabel.textContent = "Speed";

    const speedSlider = document.createElement("input");
    speedSlider.type = "range";
    speedSlider.className = "volume-slider";
    speedSlider.min = "0.5";
    speedSlider.max = "2";
    speedSlider.step = "0.1";
    speedSlider.value = settings.speed;

    const speedValue = document.createElement("span");
    speedValue.className = "setting-value";
    speedValue.textContent = settings.speed.toFixed(1) + "x";

    speedSection.appendChild(speedLabel);
    speedSection.appendChild(speedSlider);
    speedSection.appendChild(speedValue);

    // Loop time control
    const loopSection = document.createElement("div");
    loopSection.className = "sound-setting";

    const loopLabel = document.createElement("label");
    loopLabel.textContent = "Loop Wait (ms)";

    const loopSlider = document.createElement("input");
    loopSlider.type = "range";
    loopSlider.className = "volume-slider";
    loopSlider.min = "0";
    loopSlider.max = "5000";
    loopSlider.step = "100";
    loopSlider.value = settings.loopTime;

    const loopValue = document.createElement("span");
    loopValue.className = "setting-value";
    loopValue.textContent = settings.loopTime + "ms";

    loopSection.appendChild(loopLabel);
    loopSection.appendChild(loopSlider);
    loopSection.appendChild(loopValue);

    expandedView.appendChild(volumeSection);
    expandedView.appendChild(speedSection);
    expandedView.appendChild(loopSection);

    soundCard.appendChild(expandedView);
    activeSoundsList.appendChild(soundCard);

    // Event listeners
    editButton.addEventListener("click", () => {
      if (expandedSound === soundName) {
        expandedView.style.display = "none";
        compactView.style.display = "flex";
        expandedSound = null;
      } else {
        if (expandedSound) {
          // Find all cards and close the previously expanded one
          const allCards = document.querySelectorAll(".sound-card");
          allCards.forEach((card) => {
            if (card.dataset.soundName === expandedSound) {
              const expandedChild = card.querySelector(".sound-card-expanded");
              const compactChild = card.querySelector(".sound-card-compact");
              if (expandedChild) expandedChild.style.display = "none";
              if (compactChild) compactChild.style.display = "flex";
            }
          });
        }
        expandedView.style.display = "flex";
        compactView.style.display = "none";
        expandedSound = soundName;
      }
    });

    volumeSlider.addEventListener("input", (e) => {
      const volume = e.target.value / 100;
      volumeValue.textContent = Math.round(volume * 100) + "%";
      settings.volume = volume;

      const audio = audioPlayers.get(soundName);
      if (audio && !audio.paused) {
        audio.volume = globalVolume * volume;
      }
    });

    speedSlider.addEventListener("input", (e) => {
      const speed = parseFloat(e.target.value);
      speedValue.textContent = speed.toFixed(1) + "x";
      settings.speed = speed;

      const audio = audioPlayers.get(soundName);
      if (audio) {
        audio.playbackRate = speed;
      }
    });

    loopSlider.addEventListener("input", (e) => {
      const loopTime = parseInt(e.target.value);
      loopValue.textContent = loopTime + "ms";
      settings.loopTime = loopTime;
    });

    removeButton.addEventListener("click", () => {
      const audio = audioPlayers.get(soundName);
      if (audio) {
        audio.pause();
        audio.currentTime = 0;
      }
      activeSounds.delete(soundName);
      expandedSound = null;

      const button = document.querySelector(`[data-sound-name="${soundName}"]`);
      if (button) {
        button.classList.remove("active");
      }

      updateActiveSoundsDisplay();

      fetch("/api/save-sound", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ sound: soundName, action: "remove" }),
      }).catch((error) => console.error("Failed to save sound:", error));
    });
  });
}

// Initialize sound buttons
soundButtonsGrid.forEach((button) => {
  button.addEventListener("click", () => {
    const soundPath = button.dataset.soundPath;
    const soundName = button.dataset.soundName;

    initSoundSettings(soundName);

    let audio = audioPlayers.get(soundName);
    if (!audio) {
      audio = new Audio(`/public/sounds/${soundPath}`);
      audio.loop = false;
      audioPlayers.set(soundName, audio);

      audio.addEventListener("ended", () => {
        const settings = soundSettings.get(soundName);
        if (activeSounds.has(soundName) && settings) {
          setTimeout(() => {
            if (activeSounds.has(soundName)) {
              audio.currentTime = 0;
              audio.play();
            }
          }, settings.loopTime);
        }
      });
    }

    if (audio.paused) {
      const settings = soundSettings.get(soundName);
      audio.volume = globalVolume * settings.volume;
      audio.playbackRate = settings.speed;
      audio.play();
      button.classList.add("active");
      activeSounds.add(soundName);
      updateActiveSoundsDisplay();

      fetch("/api/save-sound", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ sound: soundName, action: "add" }),
      }).catch((error) => console.error("Failed to save sound:", error));
    } else {
      audio.pause();
      audio.currentTime = 0;
      button.classList.remove("active");
      activeSounds.delete(soundName);
      updateActiveSoundsDisplay();

      fetch("/api/save-sound", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ sound: soundName, action: "remove" }),
      }).catch((error) => console.error("Failed to save sound:", error));
    }
  });
});
