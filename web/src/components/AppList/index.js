import React, { useEffect, useState } from "react";
import { Swiper, SwiperSlide } from "swiper/react";
import { Navigation, Keyboard, Autoplay } from "swiper";
import Poster from "../Poster";
import { getAppList } from "../../services/api/apps";

import "swiper/css";
import "swiper/css/navigation";
import "./style.scss";

export default function AppList({ onSelectApp }) {
  const [apps, setApps] = useState([]);

  useEffect(async () => {
    const resp = await getAppList();
    if (resp.errorCode === undefined || resp.errorCode === 0) {
      if (resp.data?.apps !== null) {
        setApps(resp.data.apps);
      }
    }
  }, []);

  return (
    <div className="app-list">
      <Swiper
        spaceBetween={50}
        slidesPerView={4}
        modules={[Navigation, Keyboard, Autoplay]}
        navigation
        loop
        keyboard={{
          enabled: true,
        }}
        autoplay={{
          delay: 5000,
          disableOnInteraction: false,
        }}
      >
        {apps.map((app) => {
          return (
            <SwiperSlide
              key={app.id}
              onClick={() => {
                onSelectApp(app.id);
              }}
            >
              <Poster src={app.posterURL} />
            </SwiperSlide>
          );
        })}
      </Swiper>
    </div>
  );
}
