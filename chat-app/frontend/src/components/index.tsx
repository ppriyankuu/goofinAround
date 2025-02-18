'use client'

import { useRouter } from 'next/navigation';
import RoomSelector from '../components/room-selector';

const HomePage = () => {
  const router = useRouter();

  const handleRoomSelect = (room: string) => {
    router.push(`/${room}`);
  };

  return (
    <div className="flex items-center justify-center h-screen">
      <RoomSelector onSelect={handleRoomSelect} />
    </div>
  );
};

export default HomePage;