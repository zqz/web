export function color(k: string) : string {
  switch(k) {
    case 'application': return '#434a54';
    case 'audio': return '#fcbb42';
    case 'font': return '#37bc9b';
    case 'example': return '#da4453';
    case 'image': return '#4a89dc';
    case 'message': return '#da4453';
    case 'multipart': return '#da4453';
    case 'text': return '#967adc';
    case 'video': return '#a0d468';
  }

  return '#434a54';
}
